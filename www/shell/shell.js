// import "/node_modules/xterm/lib/xterm.js";
// import "/node_modules/xterm-addon-webgl/lib/xterm-addon-webgl.js";
// import "/node_modules/xterm-addon-fit/lib/xterm-addon-fit.js";
import { sshEntry, tkeEntry, parseTkeInit } from "./store.js";
import { Terminal } from "xterm";
import { TrzszAddon } from "trzsz";
import { WebglAddon } from "xterm-addon-webgl";
import { FitAddon } from "xterm-addon-fit";
import ShellFloat from "./shell_float.svelte";
import { hash } from "../utility.js";
import { type } from "os";

const query = new URLSearchParams(location.search);
const route = query.getAll("route");
const index = query.get("entry");
let   shell = query.get("type") || "k8s";

async function createSSH() {
    const body = {route: []};
    sshEntry.subscribe(($entry) => {
        for (const index of route) {
            const x = $entry.fetch(index);
            if (x) body.route.push(x);
        }
        const entry = $entry.fetch(index);
        document.title = `${entry.desc} (${entry.host})`;
        body.route.push(entry);
    });
    return body;
}

async function createK8S() {
    const body = {
        cluster_id: query.get("cluster_id"),
        namespace: query.get("namespace"),
        pod: query.get("pod"),
    };
    document.title = `${body.pod} (${body.cluster_id})`;
    return body;
}

async function createTKE() {
    const url = new URL("http://" + localStorage.getItem("tke:store:jump"));
    const body = {route: [ {
        host: url.hostname,
        port: parseInt(url.port),
        user: url.username,
        pass: url.password,
    } ], init: ""};
    tkeEntry.subscribe(($entry) => {
        const e = $entry.fetch(index)
        body.init = e.init;
        const entry = parseTkeInit(body.init);
        document.title = `${e.desc} (${entry.pod})`;
    });
    shell = "ssh"; // TKE 实际是 SSH 会话
    return body;
}

async function createBody() {
    switch(shell) {
    case "k8s":
        return createK8S();
    case "tke":
        return createTKE();
    default:
        return createSSH();
    }
}

async function createShell(term, body) {
    const req = JSON.stringify(body);
    const key = await hash(req, Date.now());
    const rsp = await fetch(`/bash/create?key=${key}&type=${shell}&rows=${term.rows}&cols=${term.cols}`, {
        method: "POST",
        headers: {
            "content-type": "application/json"
        },
        body: req,
    })
    const rst = await rsp.json()
    if (rst.error) throw new Error(rst.error);
    else return key;
}

async function createTerminal(key) {
    let fitting = false;
    const term = new Terminal({
        theme: {
            foreground: '#c5c8c6',
            background: '#161719',
            cursor: '#d0d0d0',

            black: '#000000',
            brightBlack: '#000000',

            red: '#fd5ff1',
            brightRed: '#fd5ff1',

            green: '#87c38a',
            brightGreen: '#94fa36',

            yellow: '#ffd7b1',
            brightYellow: '#f5ffa8',

            blue: '#85befd',
            brightBlue: '#96cbfe',

            magenta: '#b9b6fc',
            brightMagenta: '#b9b6fc',

            cyan: '#85befd',
            brightCyan: '#85befd',

            white: '#e0e0e0',
            brightWhite: '#e0e0e0'
        },
        cursorStyle: 'bar',
        // fontFamily: "Cascadia Mono",
        // fontFamily: "Intel One Mono",
        // fontFamily: "Sarasa Term SC",
        fontFamily: "Noto Sans Mono CJK SC",
        fontSize: 15,
        lineHeight: 1.2,
    });
    term.open(document.getElementById('terminal'));
    const fitAddon = new FitAddon();
    term.loadAddon(fitAddon);
    term.loadAddon(new WebglAddon());
    term.focus();
 
    let timeout;
    term.fit = function(cb) {
        clearTimeout(timeout);
        timeout = setTimeout(function() {
            fitAddon.fit();
            console.log("terminal fit: ", term.rows, "x", term.cols);
            if (cb instanceof Function) cb(term);
        }, 500);
    };
    window.term = term;
    // return new Promise((resolve) => { term.fit(resolve) });
    return term;
}

function createStream(key, term) {
    const stream = new WebSocket(`ws://127.0.0.1:8080/bash/stream?key=${key}`);
    stream.binaryType = "arraybuffer";
    // 由 TrzszAddon 接管 stream 与 Terminal 间数据交换
    const promise = new Promise((resolve) => {
        stream.addEventListener("open", function(e) {
            resolve(stream);
        });
    });
    // stream.addEventListener("message", function(e) {
    //     term.write(new Uint8Array(e.data));
    // })
    // stream.addEventListener("close", function(e) {
    //     console.log("stream close");
    //     term.write(`\r\n\u001b[0;31m${e.toString()}\u001b[0m\r\n`);
    // })
    // term.onData(function(data) {
    //     stream.send(data);
    // });
    term.loadAddon(new TrzszAddon(stream));
    term.onResize(async function(e) {
        const rsp = await fetch(`/bash/resize?key=${key}`, {
            method: "POST",
            headers: {
                "content-type": "application/json",
            },
            body: JSON.stringify({
                "rows": e.rows,
                "cols": e.cols, 
            })
        });
    });
    return promise;
}

function createKeeper(stream, term, float) {
    window.onblur = function() {
        float.$$set({"refreshing": true});
    }
    window.onfocus = function() {
        float.$$set({"refreshing": false});
    }
    let timeout, enable = false;
    const ping = function() {
        if (enable) {
            stream.send('\0');
            timeout = setTimeout(ping, 30000);
        }
    };
    term.onTitleChange(function() {
        enable = false
        clearTimeout(timeout);
    });

    float.$on("refresh", function(e) {
        if (e.detail.enable) {
            enable = true;
            setTimeout(ping, 25000);
        } else {
            enable = false;
        }
    });
    // term.onData(function() {
    //     console.log("canceld with: data", new Date());
    //     clearTimeout(timeout);
    //     timeout = setTimeout(ping, 30000);
    // });
    // term.onWriteParsed(function() {
    //     console.log("canceld with: write", new Date());
    //     clearTimeout(timeout);
    //     timeout = setTimeout(ping, 30000);
    // })
}

function createFloat(key) {
    const float = new ShellFloat({
        target: document.getElementById("float"),
        props: {
            key: key,
            fontFamily: "Intel One Mono",
        }
    });
    float.$on("font-family", function(e) {
        term.options.fontFamily = e.detail["font-family"];
        term.fit();
    });
    return float;
}

(async function() {
    const term = await createTerminal();
    const body = await createBody();
    let key;
    try {
        key = await createShell(term, body);
    }catch(e) {
        term.write(e.toString());
        return;
    }
    const float = await createFloat(key, term);
    const stream = await createStream(key, term);
    createKeeper(stream, term, float);

    window.addEventListener("resize", term.fit);
    setTimeout(term.fit, 3000);
})();