# jsonstore
JSON Database


## Common Curl Commands

```sh
curl -v -d '{"name":"foo"}' -H "Content-Type: application/json" -X POST http://localhost:9090/collections

curl -v -d '{"first":"rob"}' -H "Content-Type: application/json" http://localhost:9090/collections/foo/1

curl -v  http://localhost:9090/collections/foo/1
```
