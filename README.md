# Overview

[![Build Status](https://travis-ci.com/ctrabold/stoplight_go.svg?token=Bq7v8yEtjzfGtBmQs6m5&branch=master)](https://travis-ci.com/ctrabold/stoplight_go)

Stoplight_GO is a build monitoring tool that is largely based off [stoplight](https://github.com/customink/stoplight), and is written in [Go](https://golang.org) to make deployments easier.

Stoplight_GO has:

- built-in support for [GoCD](https://www.gocd.org/)
- built-in support for [Jenkins](https://jenkins.io/)
- built-in support for [any tool](https://github.com/erikdoe/ccmenu/wiki/Supported-Servers) that provides a `cc.xml`

Stoplight_GO is designed to be displayed on large television screens or monitors. It automatically resizes to fill the maximum real estate the screen can offer.

The main design goals was simplicity of deployment.
Feature parity with the [original Stoplight](https://github.com/customink/stoplight) was not a goal.
Stoplight_GO supports _ONE server only_ at the moment.

## Requirements

- Create a `config.yml` file. Minimal Example:

```yaml
server:
    url: http://localhost:8080/cc.xml
```

Open [config.example.yml](config.example.yml) to see all available options.

- Start a CI server that provides a `cc.xml` endpoint
  You can start a test server by running this command in a new terminal:

```bash
go run testserver/main.go
```

## Usage

- Download the binary from <https://github.com/ctrabold/stoplight_go/releases>
- Make binary executable: `chmod +x stoplight_go_*`
- Make sure you've created a `config.yml` file (see above)
- Start the binary

You can use ENV variables to override the bind address and port:

```bash
STOPLIGHT_HOST=127.0.0.1 STOPLIGHT_PORT=3000 ./stoplight_go_*
```

Default:

```bash
STOPLIGHT_HOST=0.0.0.0 STOPLIGHT_PORT=1323 ./stoplight_go_*
```

- Open a browser
- Point the browser to the URL <http://locahost:1323> or the address you've defined

## Building from Source

- Build the binary

```bash
go build stoplight.go
```

- Start the binary `./stoplight`
- Open the URL in the browser <http://0.0.0.0:1323>

For SSL support read <https://echo.labstack.com/cookbook/auto-tls>

## Cross Compiling

Execute this command:

```bash
go get github.com/mitchellh/gox
gox -os="linux darwin" -arch="amd64" -ldflags "-X main.Rev=`git rev-parse --short HEAD`" -verbose
```

TravisCI will automatically push new releases to <https://github.com/ctrabold/stoplight_go/releases>

## Known Issues & Todos

- Only one CI server is supported for now
- JavaScript and CSS assets are pre-compiled
- Tests are incomplete

## Ideas

- [ ] Implement linter <https://github.com/golangci/golangci-lint>

## Contributing

1. Fork it (<https://github.com/ctrabold/stoplight_go/fork>)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Run tests: `go test`
5. Push to the branch (`git push origin my-new-feature`)
6. Create a new Pull Request

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)
and we follow [Conventional Commits](https://www.conventionalcommits.org).

## Contributors

- [Christian Trabold](https://github.com/ctrabold) - creator and maintainer
- [Paul Yeoh](https://github.com/paul-ylz)

## Credits

- GreenScreen was original developed by [martinjandrews](https://github.com/martinjandrews/greenscreen/).
- The former version of GreenScreen was a fork of the updates made by [rsutphin](https://github.com/rsutphin/greenscreen/).
- The version of Stoplight was written by [sethvargo](https://github.com/sethvargo)
