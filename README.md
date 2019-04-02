### My extensions to [VKCOM/noverify](https://github.com/VKCOM/noverify) php linter

Adds the following checks:
- [PSR-2](https://www.php-fig.org/psr/psr-2/): "static MUST be declared after the visibility."

#### Usage

1. Clone the project to `$GOPATH/src/` and `cd noverify-psr2/`
2. Build `go build` and install `go install`
3. Optionally run tests `go test`
4. Change the path of your current `noverify` to `$GOPATH/bin/noverify-psr2` (e.g. in your [VS Code extension](https://github.com/VKCOM/noverify#visual-studio-code-integration))
