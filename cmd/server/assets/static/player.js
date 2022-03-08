'use strict';

const WASM_URL = "static/player.wasm";

var mod, inst;

function play(puzzle) {
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject)
        .then((result) => {
            mod = result.module;
            inst = result.instance;
            go.run(inst);
            fetch(puzzle)
                .then(result => result.json())
                .then((data) => {
                    loadPuzzle(data);
                })
                .catch((error) => {
                    console.warn(error);
                });
        })
        .catch((error) => {
            console.warn(error);
        });
}

function dismiss() {
    document.getElementById("instructions").style.display = "none";
}
