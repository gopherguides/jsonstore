# JSONStore

This project is not intended for production usage.  It's purpose is to show how to benchmark, profile, and optimize a Go service as part of a distributed system.

This repo is used in [Gopher Guides](https://www.gopherguides.com) advanced profiling and optimization class.

However, if you decide you want a poorly performaing JSON data store, feel free to use this project :-).


## Common Curl Commands

These are some common shorcuts to test the functionality of the key/value JSONStore api:

```sh
curl -v -d '{"name":"foo"}' -H "Content-Type: application/json" -X POST http://localhost:9090/collections

curl -v -d '{"first":"rob"}' -H "Content-Type: application/json" http://localhost:9090/collections/foo/1

curl -v  http://localhost:9090/collections/foo/1
```

Additional generic curl references for JSON data here: https://gist.github.com/subfuzion/08c5d85437d5d4f00e58
