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
        -v $PWD:/var/lib/noflake \
        -p 8000:8000 \
        ghcr.io/fact-finder/noflake:0.1.0 \
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
DTOs and interfaces for the API.

If you want to change or extend the API, edit `./openapi.yaml` first, then run
`go generate` - this invokes oapi-codegen. Now `go build` usually fails with an
error, e.g.

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

### Building the docker image

If you merge changes, Github will automatically build a new docker image. But
if you want to build the image locally, you use go build with some
extra flags to create a static binary that works in the alpine container.

```bash
$ go build -tags netgo,osusergo -ldflags="-s -w -extldflags=-static"
# github.com/FACT-Finder/noflake
/usr/bin/ld: /tmp/go-link-3319491090/000010.o: in function `unixDlOpen':
/home/thomas/src/go/pkg/mod/github.com/mattn/go-sqlite3@v1.14.13/sqlite3-binding.c:41392: warning: Using 'dlopen' in statically linked applications requires at runtime the shared libraries from the glibc version used for linking
$ ldd noflake
        not a dynamic executable
$ docker build --tag noflake:local .
Sending build context to Docker daemon  2.332GB
Step 1/6 : FROM alpine:3.17
 ---> 49176f190c7e
Step 2/6 : WORKDIR /var/lib/noflake
 ---> Using cache
 ---> 5703c60ff171
Step 3/6 : RUN mkdir /opt/noflake
 ---> Using cache
 ---> a509c9a1c229
Step 4/6 : ADD noflake /opt/noflake
 ---> Using cache
 ---> e4b0ad30903d
Step 5/6 : EXPOSE 8000
 ---> Using cache
 ---> 58dde83e08a9
Step 6/6 : ENTRYPOINT ["/opt/noflake/noflake"]
 ---> Using cache
 ---> 76ad7efec763
Successfully built 76ad7efec763
Successfully tagged noflake:local
```

### Releasing a new version

Please create a new release [via GitHub](https://github.com/FACT-Finder/noflake/releases/new) on a new tag. 
This will trigger our release workflow which pushes the newest image to the docker registry.