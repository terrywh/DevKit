
import { spawn } from "node:child_process"

async function build(name, os, arch) {
    // console.log(`\x1b[38;5;214m${name}\x1b[0m_\x1b[38;5;204m${os}\x1b[0m_\x1b[38;5;194m${arch}\x1b[0m`)
    console.time(`${name}_${os}_${arch}`)
    return new Promise((resolve, reject) => {
        const cp = spawn("go", [
            "build",
            "-o", `./bin/${name}_${os}_${arch}`,
            `./cmd/${name}`,
        ],{
            stdio: "inherit",
            env: Object.assign({}, process.env, {
                "GOOS": os,
                "GOARCH": arch,
            }),
        });
        cp.on("close", (code) => {
            console.timeEnd(`${name}_${os}_${arch}`)
            code == 0 ? resolve() : reject(code);
        });
    });
}
// 链接使用 Webview 时最好在对应平台编译
await build("devkit-client", "darwin")
await build("devkit-client", "darwin")

await build("devkit", "windows", "amd64")
await build("devkit","linux","amd64")
await build("devkit","darwin","arm64")

await build("devkit-server","windows","amd64")
await build("devkit-server","linux","amd64")
await build("devkit-server","darwin","arm64")

await build("devkit-relay","linux","amd64")
