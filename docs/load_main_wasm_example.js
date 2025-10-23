const go = new Go(); // defined in wasm_exec.js

WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject)
    .then((result) => {
        go.run(result.instance);
        let script = document.createElement("script")
        script.textContent = example()
        document.body.appendChild(script)
    });