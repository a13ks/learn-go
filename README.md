# Demo project to learn Go

## Adding new album to the albums database

```
curl http://localhost:8080/albums \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'
```

## Useful links

- https://go.dev/doc/tutorial/web-service-gin
- https://betterprogramming.pub/build-a-scalable-api-in-go-with-gin-131af7f780c0
