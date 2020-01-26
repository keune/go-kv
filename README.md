# go-kv
This is a simple REST API that works as a key-value service. It's written to demonstrate a simple dockerized Go Web application. It uses Redis for storage.

## Requirements
- Install docker && docker-compose

## Run this application
Copy .env.example to .env and change values as you need.
Run `docker-compose up`

## Usage
To create a new key, make a POST request to /new/yourKey. The application should return 200 with a new URL that you can use to set value.

```console
# curl -X POST localhost:8090/new/testKey
http://localhost:8090/AvuY9sQ/testKey
```

To set value to a key, append your value to the URL returned from /new and make a POST request. The application should return a 200 response with an empty body.

```console
# curl -X POST http://localhost:8090/AvuY9sQ/testKey/testValue
```

To get value, simply make a GET request to the key's URL. The application should return a 200 response with the key's value in the body.
```console
# curl http://localhost:8090/AvuY9sQ/testKey
testValue
```
