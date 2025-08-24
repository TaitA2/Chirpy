# Chirpy

* A Twitter-like backend server built in Go following a guide on [Boot.Dev](https://www.boot.dev/u/taita2) making use of JWTs, refresh tokens, webhooks, and API keys.

## Features
* Create and login to user accounts using email and password authentication.
* Create, edit, view, and delete text-based "chirps" stored in a PostgreSQL database using JWTs.
* Track premium account (Red) status using API Key verified webhooks.
* Generate new access tokens using a refresh token.
* Revoke access to refresh tokens.
* View site usage metrics

## Requirements
* Go 1.22+
* PostgreSQL 15

## Installation
* Chirpy can be installed using the following command:
* `go install github.com/TaitA2/Chirpy@latest`

## Usage
### Users
* Users can be created or logged into with a POST request to the /api/users endpoint in the following format:
  * `REQUEST: {"email": "name@example.com", "password": "password123"}`
  * `RESPONSE: {"ID": "UUID", "CreatedAt": "datetime", "UpdatedAt": "datetime", "email": "name@example.com", "Token": "Access-token", "RefreshToken": "Refresh-token", "IsChirpyRed": false}`
* User credentials can be updated with a PUT request to the /api/users endpoint in the following format:
  * `REQUEST: {"email": "newemail@example.com", "password": "newpassword123"} HEADER: "Authorization: Bearer <Access-Token>"`
  * `RESPONSE: {"ID": "UUID", "CreatedAt": "datetime", "UpdatedAt": "datetime", "email": "name@example.com", "Token": "Access-token", "RefreshToken": "Refresh-token", "IsChirpyRed": false}`
* Users can generate a new access token with a POST request to the /api/refresh endpoint in the following format:
  * `REQUEST: {} HEADER: "Authorization: Bearer <Refresh-Token>"`
  * `RESPONSE: {"token": "<Access-Token>"}`
* Refresh tokens can be revoked with a POST request to the /api/revoke endpoint in the following format:
  * `REQUEST: {} HEADER: "Authorization: Bearer <Refresh-Token>"`
  * `RESPONSE: 204`
* NOTE
  * Access tokens expire after 1 hour
  * Refresh tokens expire after 60 days.

### Chirps
* Chirps can be created with a POST request to the /api/chirps endpoint in the following format:
  * `REQUEST: {"body": "body of the chirp"} HEADER: "Authorization: Bearer <Access-Token>"`
  * `RESPONSE: {"ID":"<Chirp-ID>","CreatedAt":"<DateTime>","UpdatedAt":"<DateTime>","Body":"<body of the chirp>","UserID":"<User-UUID>"}`
* Chirps can be deleted with a DELETE request to the /api/chirps/{chirp-id} endpoint 
  * `REQUEST: {} HEADER: "Authorization: Bearer <Access-Token>"`
  * `RESPONSE: 204`
