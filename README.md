# HHot

## Install

```shell
go install github.com/SamHennessy/hhot/cmd/hhot@latest
```

## Tips

### Add Go bin folder to your $PATH

```shell
export PATH=$(go env GOPATH)/bin:$PATH
```

## Development

Build and install locally

```shell
go install ./cmd/hhot
```

## TODO

- Split the hot reload tool and framework (HFW)?

### Asset Manager

- Embed all files
- On start scan all files and create a hash
  - Maybe do in a goroutine
  - So reading may not always work
- Have a function that will add a hash to the image URL
- Allow adding of CSS and JS that is not app.js and app.css
- Add hash to fav icons
  - Allow optional template variables to the manifest file for a hash?
  - Or just use a manuel version number?
- Images
  - Help to get image dimensions to prevent content shift
  - Build an img or picture tag?
    - What about CSS backgrounds?
  - Alternative encoding webp?
  - Downscale to match dpi?
- Native GZip and Brotli support?
- ETags

### Router

- When navigating to a new page, scroll to top of the page

### Hot Reload

- Move to a https://wails.io/ app?

- Store logs in SQLite
- Allow view of logs made before the ui was loaded
- filter logs
- resize log window
- hide log window
- Make address bar interactive