# Razorblade URL shortener

This is a small educational project that shows how easy it is to run an http server using Go.

## Some features

This shortener has several endpoints:
- **/** - there is a shortening form 
- **/ping** - checks PostgreSQL database connection
- **/api/user/urls** - returns all links shortened by the current user
- **/api/shorten** - accepts a link to shorten in JSON format
- **/api/shorten/batch** - accepts a batch of links to shorten in JSON format
- **DELETE /api/user/urls** - accepts a batch of links to delete also in JSON format

The gRPC API is also available.

## How to use 

Just run /cmd/shortener/main.go

It is recommended to have a PostgreSQL database. Just specify its address in internal/app/config/config.go, the required table will be created automatically. However, it is possible to work without it. In this case, the storage.txt file will be used.