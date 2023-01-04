# Noflake

A web service that automatically detects and tracks flaky tests based on junit
reports.

## How it works

After every CI build of your software, you upload all junit reports together
with the commit SHA of the test run. Noflake keeps track of the results and
automatically detects flakiness when a test passes and fails on the same
commit. The Noflake UI shows you all detected flaky tests in a given time range
sorted by the number of fails relative to the total test invocations.

## How to use

1. Run Noflake:
   ```bash
    docker run \
        -v $PWD/noflake.sqlite3:/opt/noflake/noflake.sqlite3 \
        -p 8000:8000 \
        ghcr.io/fact-finder/noflake:0.0.9 \
        --token verysecuretoken
   ```
2. Upload junit reports from your CI server after every test execution. For a
   GitLab job it might look like this:
   ```bash
   curl -H "Authorization: Bearer verysecuretoken"
   $(find -path '*/build/test-results/**/TEST-*.xml' | sed -E 's/^(.*)$/-F report=@"\1"/')
   "http://your-noflake-instance.com:8000/report/${CI_COMMIT_SHA}?url=${CI_JOB_URL}"
   ```
3. Open the Noflake UI and find the most annoying flaky tests:
   http://your-noflake-instance.com:8000/ui

## Development

### Changing the API

Noflake is written in Go with an API first approach and SQLite3 as database.

It uses [oapi-codegen](https://github.com/deepmap/oapi-codegen) to generate
DTOs and interfaces for the API. You can install it with

```bash
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.11.0
```

and use it with `go generate`. Make sure that your `$PATH` includes the go
binary directory which is either `$(go env GOBIN)` or `$(go env GOPATH)/bin`
(usually `~/go/bin`).

If you want to change or extend the API, edit `./openapi.yaml` first, then run
`go generate`. Now `go build` usually fails with an error, e.g.

```
$ go build
# github.com/FACT-Finder/noflake/api
api/api.go:37:45: cannot use api (variable of type *api) as type ServerInterface in struct literal:
        *api does not implement ServerInterface (missing GetHello method)
```

The error indicates that the `api` struct doesn't implement `ServerInterface`
which was generated from the openapi definition because it misses a
`GetHello` method. In this case you would need to define the method
`GetHello` on the `api` struct.

In addition you have to register the generated method `wrapper.GetHello`
with a corresponding route in `api/api.go`:

```go
app.GET("/hello", wrapper.GetHello)
```

Make sure that the route in `api/api.go` matches the openapi definition because
the compiler won't check it.

### Changing the database

If you want to change the database structure, add new up and down migrations to
the `database/migrations` folder. Update `model/model.go` and all usages of the
models accordingly. The not yet applied migrations are automatically executed
on app startup. You can also manually migrate up and down using the
[golang-migrate cli](https://github.com/golang-migrate/migrate#cli-usage).
