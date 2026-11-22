# gui

The console window shared by the trimlet apps: a dark text console driven only by the keyboard, drawn
with [Fyne](https://fyne.io), with the trimlet color palettes, the colored title bar on macOS, a
JSON/XML viewer and the "Quit …?" window.

It is not an app. The apps ([db](https://github.com/trimlet/db), [endpoint](https://github.com/trimlet/endpoint),
[files](https://github.com/trimlet/files), [seal](https://github.com/trimlet/seal),
[tools](https://github.com/trimlet/tools), [auth](https://github.com/trimlet/auth),
[password](https://github.com/trimlet/password)) import it, and `go build` fetches it by itself:

```go
import "github.com/trimlet/gui"
```

## How it is made

The code is written with Claude Opus (Anthropic), under the direction and review of a developer with
almost 30 years of experience. Every feature starts from a real, repeated need, and most ideas are
turned down: the goal is to stay small.

## License

GPL-3.0, see [LICENSE](LICENSE).
