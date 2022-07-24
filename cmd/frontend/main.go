package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
)

type accessToken struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	Type         string `json:"token_type"`
}

const tokenURL = "http://localhost:8080/auth/realms/Customers/protocol/openid-connect/token"
const authURL = "http://localhost:8080/auth/realms/Customers/protocol/openid-connect/auth"

func main() {

	/**
		Este par de verifier/challenge deve ser criado por request, isso é só um teste, no frontend pode ser salvo na Session Storage.
	**/
	codeVerifier, _ := cv.CreateCodeVerifier()
	codeChallenge := codeVerifier.CodeChallengeS256()

	// 1) Redireciona o Usuário para o Keycloak para fazer o login
	http.HandleFunc("/auth", handleAuth(codeChallenge))

	// 2) Recebe o Auth Code e troca por um Access Token
	http.HandleFunc("/callback", handleCallback(codeVerifier))

	log.Fatal(http.ListenAndServe(":3030", nil))

}

func handleAuth(codeChallenge string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, _ := url.Parse(authURL)
		q := u.Query()
		q.Set("client_id", "test-client-public-spa")            // Required
		q.Set("response_type", "code")                          // Required
		q.Set("redirect_uri", "http://localhost:3030/callback") // Optional
		q.Set("state", "mystate")                               // Recommended
		q.Set("scope", "openid")                                // Required
		q.Set("code_challenge_method", "S256")                  // Required for PKCE
		q.Set("code_challenge", codeChallenge)                  // Required for PKCE - [HOW GENERATE CODE VERIFIER] https://datatracker.ietf.org/doc/html/rfc7636#page-8

		u.RawQuery = q.Encode()
		fmt.Println(u)
		http.Redirect(w, r, u.String(), http.StatusSeeOther)
	}
}

func handleCallback(codeVerifier *cv.CodeVerifier) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query()["code"][0]
		token := getAccessToken(code, codeVerifier.String())

		fmt.Println()
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(token)
	}
}

func getAccessToken(code, codeVerifier string) *accessToken {
	data := url.Values{}
	data.Set("client_id", "test-client-public-spa")
	data.Set("code_verifier", codeVerifier)
	data.Set("redirect_uri", "http://localhost:3030/redirect")
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)

	fmt.Println(data.Encode())

	r, err := http.NewRequest(http.MethodPost, tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		s, _ := io.ReadAll(resp.Body)
		fmt.Println(string(s))
		return nil
	}

	var token accessToken
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		panic(err)
	}
	return &token
}
