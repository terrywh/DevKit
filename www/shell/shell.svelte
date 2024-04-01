
<script>
// import "/node_modules/xterm/lib/xterm.js";
// import "/node_modules/xterm-addon-webgl/lib/xterm-addon-webgl.js";
// import "/node_modules/xterm-addon-fit/lib/xterm-addon-fit.js";
import { onMount } from "svelte";
import { sshEntry, tkeEntry, parseTkeInit } from "./store.js";
import ShellFloat from "./shell_float.svelte";
import ShellTerminal from "./shell_terminal.svelte";

import { TrzszAddon } from "/node_modules/trzsz/lib/trzsz.mjs";
// import { TrzszAddon } from "https://esm.sh/trzsz"; // 内部存在不支持的模块
import { hash } from "../utility.js";

const query = new URLSearchParams(location.search);
const route = query.getAll("route");
const index = query.get("entry");
let   shell = query.get("type") || "k8s";

async function createSSH() {
    const body = {route: []};
    sshEntry.subscribe((store) => {
        for (const index of route) {
            const x = store.fetch(index);
            if (x) body.route.push(x);
        }
        const entry = store.fetch(index);
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
    tkeEntry.subscribe((store) => {
        const e = store.fetch(index)
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

async function createShell(body) {
    const req = JSON.stringify(body);
    const key = await hash(req, Date.now());
    const rsp = await fetch(`/bash/create?key=${key}&type=${shell}&rows=60&cols=90`, {
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

let cTerminal;
let refreshing = $state(false);

function createStream(key) {
    const stream = new WebSocket(`ws://127.0.0.1:8080/bash/stream?key=${key}`);
    stream.binaryType = "arraybuffer";
    // 由 TrzszAddon 接管 stream 与 Terminal 间数据交换
    const promise = new Promise((resolve) => {
        stream.addEventListener("open", function(e) {
            // cTerminal.loadAddon(new TrzszAddon(stream));
            resolve(stream);
        }, {once: true});
    });
    return promise;
}

// function createKeeper(stream, term, float) {
    
//     let timeout, enable = false;
//     const ping = function() {
//         if (enable) {
//             stream.send('\0');
//             timeout = setTimeout(ping, 30000);
//         }
//     };
//     term.onTitleChange(function() {
//         enable = false
//         clearTimeout(timeout);
//     });

//     // float.$on("refresh", function(e) {
//     //     if (e.detail.enable) {
//     //         enable = true;
//     //         setTimeout(ping, 25000);
//     //     } else {
//     //         enable = false;
//     //     }
//     // });
//     // term.onData(function() {
//     //     console.log("canceld with: data", new Date());
//     //     clearTimeout(timeout);
//     //     timeout = setTimeout(ping, 30000);
//     // });
//     // term.onWriteParsed(function() {
//     //     console.log("canceld with: write", new Date());
//     //     clearTimeout(timeout);
//     //     timeout = setTimeout(ping, 30000);
//     // })
// }

let key = $state(""), stream;

onMount(async function() {
    const body = await createBody();
    console.log("create shell: ", body.route);
    try {
        key = await createShell(body);
    }catch(e) {
        cTerminal.write(e.toString());
        return;
    }
    stream = await createStream(key);
    stream.addEventListener("close", function(e) {
        // const colors = []
        // for (let i = 16; i < 256; i++) {
        //     colors.push(`\x1b[38;5;${i}m${i}\x1b[0m`)    
        //     if ((i -4) % 12 === 11) {
        //         colors.push('\n')
        //     }
        // }
        cTerminal.write(`\r\n\x1b[38;5;202mConnection Closed, closing windows in 5s ...\x1b[0m\r\n`);
        setTimeout(() => window.close(), 4500);
    })
    // createKeeper(stream, term, float);
    
    window.addEventListener("resize", function() {
        cTerminal.fit();
    });
});

let refreshingTimeout;
function refresh() {
    clearTimeout(refreshingTimeout);
    refreshingTimeout = setTimeout(function () {
        cTerminal.fit();
        if (refreshing) refresh();
    }, 10000);
}

$effect(() => {
    if (refreshing) refresh();
    else clearTimeout(refreshingTimeout);
});

</script>

<div style="width: 100%; height: 100%;">
    <ShellTerminal bind:this={cTerminal} key={key} />
</div>
<div style="position: fixed; top: 2rem; right: 4rem; z-index: 100;">
    <ShellFloat bind:refreshing={refreshing} key={key} />
</div>
