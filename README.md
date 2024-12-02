## Technical Stack
Go: Programming language.
Echo: Web framework for handling HTTP requests.
GORM: ORM library for Go.
PostgreSQL with PostGIS: Database system with support for geographic objects.
Docker and DockerCompose: Containerization of the application and database

## Environment Variables

Set the following environment variables in your .env file or export them in your shell:

- DB_HOST: Database host
- DB_USER: Database user
- DB_PASS: Database password
- DB_NAME: Database name
- DB_PORT: Database por.
- JWT_SECRET_KEY: Secret key for JWT
- JWT_DURATION: Duration for JWT expiration
- LOG_ENV: Logging environment (development or production)

## API Endpoints

The server will start on port 8080, and you can access the API at http://localhost:8080/ , here are some CURL commands to interact with the API:

* remember the valid value for gender is MALE or FEMALE

```
# Create a Random User
curl -X POST http://localhost:8080/user/create

# Login
curl -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{
        "email": "mariliemetz@moore.com",
        "password": "WidMJHF_q51?"
    }'

# Discover Matches
curl -X GET "http://localhost:8080/discover?lat=34.0522&lng=-118.2437&distance=10000&gender=MALE&minAge=18&maxAge=50" \
    -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Swipe on a Profile
curl -X POST http://localhost:8080/swipe \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer YOUR_JWT_TOKEN" \
    -d '{
        "targetUserId": 2,
        "preference": "YES"
    }'

```

## Example of Match Curls
```

❯ curl -X POST http://localhost:8080/user/create
{"result":{"age":64,"email":"tadward@abernathy.name","gender":"FEMALE","id":6,"name":"Tamara Miller","password":"7sT32u#9UR45"}}
❯ curl -X POST http://localhost:8080/user/create
{"result":{"age":22,"email":"maryryan@kreiger.net","gender":"FEMALE","id":7,"name":"Sid Sauer","password":"hj.xBtDuf!yd"}}

❯ curl -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{
        "email": "tadward@abernathy.name",
        "password": "7sT32u#9UR45"
    }'
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTUwNzE1NjIsIlVzZXJJRCI6Nn0.VtD8pBZwMdCdZuutp31Mv1FQ-lhGB5gyPB0yMb2l-oQ"}
❯     curl -X POST http://localhost:8080/login \
    -H "Content-Type: application/json" \
    -d '{
        "email": "maryryan@kreiger.net",
        "password": "hj.xBtDuf!yd"
    }'
{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTUwNzE1NzIsIlVzZXJJRCI6N30.mxEtEHxvJ5UqSDzNCpiDxAeS-sxoE3E37jgTKPu0vR8"}


❯ curl -X POST http://localhost:8080/swipe \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTUwNzE1NjIsIlVzZXJJRCI6Nn0.VtD8pBZwMdCdZuutp31Mv1FQ-lhGB5gyPB0yMb2l-oQ" \
    -d '{
        "targetUserId": 7,
        "preference": "YES"
    }'
{"results":{"matched":false}}
❯   curl -X POST http://localhost:8080/swipe \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTUwNzE1NzIsIlVzZXJJRCI6N30.mxEtEHxvJ5UqSDzNCpiDxAeS-sxoE3E37jgTKPu0vR8" \
    -d '{
        "targetUserId": 6,
        "preference": "YES"
    }'
{"results":{"matched":true,"matchID":1}}


```
