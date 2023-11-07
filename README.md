# CSI
> Client Side Include

## Introduction

As opposed to server side includes (SSI); this tooling aims to allow the include of responses (e.g. remote HTML document
partials) in
local ones. It might be generally useful for static site generation when used in conjunction with a headless/hybrid CMS.

## Get Started

Since it might be generally useful for `HTML` generation here an example of the usage in `HTML`, but in general the
interpreter reads any document regardless of extension. And includes whatever the `HTTP GET` request returns.

```html

<html>
    <body>
        <!--{GET http://example.com}--> 
    </body>
</html>
```

The `<!--{GET http://example.com}-->` will be replaced with whatever a `HTTP GET` request will fetch from the following
host.  

Error occurring during request or URL parsing will be output as log messages and the content will be replaced by an
empty string.

## Tooling Binaries

### build

```bash
go run ./cmd/build/ --help
Usage of ...
  -file string
    	the file to convert
```

This will process the input file and output the content, with processing of the CSI snippets to stdout.

### watch

```bash
go run ./cmd/watch/ --help
Usage of ...
  -file string
    	the file to convert
  -out string
    	the target to output (default "out.html")
```

This will process the input file and will on every change; process and store the processed file as defined in `out`.

### serve

```bash
go run ./cmd/serve/ --help
Usage of ...
  -file string
    	the file to convert
  -port string
    	the target port to listen (default "4321")
```

This will process the file (similar to watch) and serve it on `127.0.0.1:*port*`.

## Building

```bash
make build # will create all binaries in a ./dist directory
```

## License

MIT
