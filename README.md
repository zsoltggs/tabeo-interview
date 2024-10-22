# Bookings

## How to run the application

```
make build-docker
docker-compose up
```


### Testing

`make test` runs all unit tests.

`make verify` runs the linter, and runs all unit, integration and E2E tests (make sure you run docker-compose up before verify).

## How to make better

- Upgrade cache, use time based expiration and then the time based decision can be dropped
- Launch date should be validated during creation
- It would be better if list bookings performs a text search (e.g. the query endpoint for spacex)
- Change E2E tests go gingko tests that are easy to read 
- I think it is also nice to have a get by id endpoint, but out of scope