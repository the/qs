# qs

A command line URL query string parser.

## Installation

If you have go installed you can easily install `qs` by running `go get -u github.com/the/qs`. Ensure you have `$GOPATH/bin` in your shell's `PATH` variable or execute it with `$GOPATH/bin/qs`.

## Usage

By default `qs` reads and parses URLs from `stdin`. You can specify to read from a file with `-f filename`, e.g. `qs -f urls.txt`.

By default `qs` simply parses the URLs and outputs them. If you specify query string parameter names those parameters will get highlighted. For example `echo "https://www.example.com/path/?name=value&foo=bar&param=value" | qs name param` will highlight the parameter names and values of `name` and `param`.

If you want to do some enhanced URL or query parameter processing you can enable JSON output with `-json` and then pipe it through [jq](https://stedolan.github.io/jq/) to do further processing. In JSON mode you can filter out query parameters by appending a list of query parameter names, otherwise all query parameters are added. Example output from `echo "https://www.example.com/path/?name=value&foo=bar&param=value" | qs -json name param | jq '.'`:

```javascript
{
  "url": "https://www.example.com/path/?name=value&foo=bar&param=value",
  "scheme": "https",
  "host": "www.example.com",
  "path": "/path/",
  "query": {
    "name": "value",
    "param": "value"
  }
}
```
