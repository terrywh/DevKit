<script>
    import { onMount } from "svelte";
    import { Terminal } from "https://esm.sh/xterm@5.3.0";
    import { WebglAddon } from "https://esm.sh/xterm-addon-webgl@0.16.0";
    import { FitAddon } from "https://esm.sh/xterm-addon-fit@0.8.0";

    let { key = "" } = $props();
    /** @type {Terminal} */
    let iTerminal;
    let iFitAddon;
    onMount(async function () {
        iTerminal = new Terminal({
            theme: {
                foreground: "#c5c8c6",
                background: "#161719",
                cursor: "#d0d0d0",

                black: "#000000",
                brightBlack: "#000000",

                red: "#fd5ff1",
                brightRed: "#fd5ff1",

                green: "#87c38a",
                brightGreen: "#94fa36",

                yellow: "#ffd7b1",
                brightYellow: "#f5ffa8",

                blue: "#85befd",
                brightBlue: "#96cbfe",

                magenta: "#b9b6fc",
                brightMagenta: "#b9b6fc",

                cyan: "#85befd",
                brightCyan: "#85befd",

                white: "#e0e0e0",
                brightWhite: "#e0e0e0",
            },
            cursorStyle: "bar",
            // fontFamily: "Cascadia Mono",
            fontFamily: "Intel One Mono",
            // fontFamily: "Sarasa Term SC",
            // fontFamily: "Noto Sans Mono CJK SC",
            fontSize: 16,
            lineHeight: 1.2,
        });
        iTerminal.open(document.getElementById("terminal"));
        iFitAddon = new FitAddon();
        iTerminal.loadAddon(iFitAddon);
        iTerminal.loadAddon(new (WebglAddon.WebglAddon || WebglAddon)());
        iTerminal.focus();
        iTerminal.onResize(async function (e) {
            const rsp = await fetch(`/bash/resize?key=${key}`, {
                method: "POST",
                headers: {
                    "content-type": "application/json",
                },
                body: JSON.stringify({
                    rows: e.rows,
                    cols: e.cols,
                }),
            });
        });

        const done = iTerminal.onTitleChange(function (e) {
            fit();
            done.dispose();
        });
    });

    let fitTimeout;
    export function fit(cb) {
        clearTimeout(fitTimeout);
        fitTimeout = setTimeout(function () {
            iFitAddon.fit();
            console.log("terminal fit: ", iTerminal.rows, "x", iTerminal.cols);
            if (cb instanceof Function) cb();
        }, 300);
    }

    export function loadAddon(addon) {
        iTerminal.loadAddon(addon);
    }

    export function write(message) {
        iTerminal.write(message);
    }
</script>

<div id="terminal" style="width: 100%; height: 100%;"></div>
