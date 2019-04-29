# cacher

I'm trying out gRPC and have made a mem-cache and a small performance test.

Get it:
```
git clone https://github.com/robinbraemer/cacher.git && cd cacher
```

Run the server:
```
go run server/*
```

Run the client:
```
go run client/*
```

**Commands:**
- set [key] [string] - Sets an entry
- get [key] - Gets an entry
- del [key] - Deletes an entry
- all - Gets all entries using server-side stream
- empty - Empty the cache
- clear - Clear console
- fill [count] - Fills the cache and benefits from HTTP/2 connection reuse
- slow-fill [count] - Fills the cache without connection reuse

**Results on my PC**
- `fill-slow 10000` takes `10.481s` to `14.674s`
- `fill 10000` takes `2s` to `3.341s`