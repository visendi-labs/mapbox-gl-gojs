# Notes on docs
- Built using https://docsify.js.org/#/.
- In the repo locally, from `/docs` folder, run: `docsify serve`. How to install docsify: https://docsify.js.org/#/quickstart.
    - You'll need to WASM-compile all the Go code (`common.go` in the example folders) first. Run `./wasm_compile_all.sh` inside `/docs`. It will output `main.wasm` files in the example folders, which will be picked up by the docsify HTML example files. See more below on WASM.
        - Before you do this you need to export a Mapbox access token, to be picked up by the Go/WASM code. Run `export MAPBOX_ACCESS_TOKEN=pk.eyJ1I...`.
- There's a bunch of examples to try out - not only by running locally but also live "serverside" examples running on the docs page. 
    - Literal Go code run in your browser, it's compiled to WASM.
	- Some serverside interactivity examples are also running actual, but mocked, HTMX (https://htmx.org) calls to the WASM browser Go "fake" backend. Essentially the same thing as if you spun up the example webserver+web client locally. 
    - Go comes with great WASM support. The file `wasm_exec.js`, which contain tools to interact with your WASM compiled Go code from JS, is taken straight from Go (shipped with Go as `$(go env GOROOT)/lib/wasm/wasm_exec.js`). Copyright 2018 The Go Authors. All rights reserved. https://go.dev/wiki/WebAssembly.
    - Some adaptations has been made to the Go examples in order to make them runnable in a WASM JS context. E.g. Go code that open files can't be WASM-ed, so `//go:embed icon.png` has been used instead.
    - Some adaptations has been made to the Mapbox-GL-GOJS examples in order for the "fake" WASM backend "responses" to be able to be executed by JS in the docs page. The browser won't execute a JS HTML string with content put straight in the DOM (added like `.innerHTML += content`). However, it will execute JS code if put into a `document.createElement("script")`. Mapbox-GL-GOJS supports wrapping everything in `<script>` tags (`mapboxglgojs.NewScript(...)`, use it!) but that didn't really work in the docs web examples (without resorting to like `eval` or some string parsing).


