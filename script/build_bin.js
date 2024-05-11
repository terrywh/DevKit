
import { spawn } from "node:child_process"

async function build(name, os, arch, ext) {
    // console.log(`\x1b[38;5;214m${name}\x1b[0m_\x1b[38;5;204m${os}\x1b[0m_\x1b[38;5;194m${arch}\x1b[0m`)
    if (os == "windows") {
        ext = ".exe"
    } else {
        ext = ""
    }
    const label = `build ${name}_${os}_${arch}${ext}`;
    console.time(label);
    return new Promise((resolve, reject) => {
        const cp = spawn("go", [
            "build",
            "-o", `./bin/${name}_${os}_${arch}${ext}`,
            `./cmd/${name}`,
        ],{
            stdio: "inherit",
            env: Object.assign({}, process.env, {
                "GOOS": os,
                "GOARCH": arch,
            }),
        });
        cp.on("close", (code) => {
            console.timeEnd(label);
            code == 0 ? resolve() : reject(code);
        });
    });
}

await build("devkit-client", "darwin", "arm64")
await build("devkit-client", "windows", "amd64")

await build("devkit", "windows", "amd64")
await build("devkit","linux","amd64")
await build("devkit","darwin","arm64")

await build("devkit-server","windows","amd64")
await build("devkit-server","linux","amd64")
await build("devkit-server","darwin","arm64")

await build("devkit-relay","linux","amd64")
