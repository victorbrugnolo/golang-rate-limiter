# Golang rate limiter

A way to protect your sever of attacks.

## Rate limiter rules
- This rate limiter allows you to use two configuration methods: by IP or by token. 
- If you send a token in the request, this token must be previously configured, otherwise you will receive an error response in the request.
- Per-token settings will override IP settings.
- Token must send in the request header on format `API_KEY=tokenname`

## Configuration
All configurations must be made in the `.env` file in the project root, following this pattern:

### By IP
```
RATE_LIMITER_BY_IP_LIMIT=5
RATE_LIMITER_BY_IP_WINDOW=10
RATE_LIMITER_BY_IP_BLOCK_WINDOW=5
```

### By token:
```
RATE_LIMITER_BY_TOKEN_tokenname_LIMIT=3
RATE_LIMITER_BY_TOKEN_tokenname_WINDOW=5
RATE_LIMITER_BY_TOKEN_tokenname_BLOCK_WINDOW=10
```

_**You can determine the name of the token and configure as many as you want**_

## Running
```bash
go mod tidy

docker-compose up -d
```

The server will be running on port `:8080`

## Tests
### Automated tests:
```bash
go test ./...   
```

### Test rate limiter by ip:
```bash
bash bash/test_by_ip.bash   
```

### Test rate limiter by token:
```bash
bash bash/test_by_token.bash 
```