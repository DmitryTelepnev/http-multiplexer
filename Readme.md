# Http Multiplexer

see [conditions](Conditions.md)

## Requirements
* docker

## Test it
```bash
make test
```

## Run it

```bash
docker build -t http-multiplexer-dtelepnev .
docker run -p 8080:8080 http-multiplexer-dtelepnev
```

## Try it

```
POST http://localhost:8080/
Accept: */*
Cache-Control: no-cache
Content-Type: application/json

[
  "https://jsonplaceholder.typicode.com/todos/1",
  "https://jsonplaceholder.typicode.com/posts/1",
  "https://jsonplaceholder.typicode.com/posts",
  "https://jsonplaceholder.typicode.com/comments",
  "http://localhost"
]


POST http://localhost:8080

HTTP/1.1 500 Internal Server Error
Content-Type: application/json
Date: Wed, 23 Dec 2020 07:49:33 GMT
Content-Length: 97

{
  "msg": "Get \"http://localhost\": dial tcp 127.0.0.1:80: connect: connection refused",
  "data": {}
}
```

```
POST http://localhost:8080/
Accept: */*
Cache-Control: no-cache
Content-Type: application/json

[
  "https://jsonplaceholder.typicode.com/todos/1",
  "https://jsonplaceholder.typicode.com/posts/1",
  "https://jsonplaceholder.typicode.com/posts",
  "https://jsonplaceholder.typicode.com/comments"
]


POST http://localhost:8080

HTTP/1.1 200 OK
Content-Type: application/json
Date: Wed, 23 Dec 2020 07:50:07 GMT
Transfer-Encoding: chunked

{
  "msg": "Ok",
  "data": {
    "https://jsonplaceholder.typicode.com/comments": "[\n  {\n    \"postId\": 1,\n    \"id\": 1,\n    \"name\": \"id labore ex et quam laborum\",\n    \"email\": \"Eliseo@gardner.biz\",\n    \"body\": \"laudantium enim quasi est quidem magnam voluptate ipsam eos\\ntempora quo necessitatibus\\ndolor quam autem quasi\\nreiciendis et nam sapiente accusantium\"\n  },\n  {\n    \"postId\": 1,\n    \"id\": 2,\n    \"name\": \"quo vero reiciendis velit similique earum\",\n    \"email\": \"Jayne_Kuhic@sydney.com\",\n    \"body\": ......

```