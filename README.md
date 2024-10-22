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
- Add a proper mock in docker compose for the third-party dependency. Right now always the real one is called.

### Assumption

> Every day you change the destination for all the launchpads. Every day of the week from the same launchpad has to be a “flight” to a different place.

My assumption that this requirement is not needed in the actual implementation. 