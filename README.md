# Keycloak + PKCE example

## How to create a PKCE client on Keycloack
1. Create a client
2. Set `Access Type` to `public`
3. Set `Standard Flow Enabled` to `ON`
4. Set `Valid Redirect URIs` to URL of your application, for example `http://localhost:3030/*`, here you can use wildcard, this address must be the same that the application will use as callback after user signin on keycloak login page.
5. Set `Base URL` for some public accesss of your application, it will be used in case of to cancel authentication then Keycloak will recirect user for this URL
6. In `Advanced Settings > Proof Key for Code Exchange Code Challenge Method` select `S256`
7. Save the configuration!