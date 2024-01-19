

function architecture() {
    let arch = "amd64";
    if (Deno.build.arch == "aarch64") {
        arch = "arm64";
    }
    return arch;
}

function executable() {
    const arch = architecture();
    if (Deno.build.os == "windows") {
        return `bin/kubectl_${Deno.build.os}_${arch}.exe`;
    } else {
        return `bin/kubectl_${Deno.build.os}_${arch}`;
    }
}

async function latest() {
    const rsp = await fetch("https://dl.k8s.io/release/stable.txt");
    return await rsp.text();
}

async function already(version) {
    const name = executable();
    try {
        await Deno.stat(name)
    } catch(ex) {
        return false;
    }
    const cmd = new Deno.Command(name, {"args": ["version", "-o", "json"]});
    const out = await cmd.output();
    const txt = new TextDecoder().decode(out.stdout);
    const ver = JSON.parse(txt).clientVersion.gitVersion;
    return ver == version;
}

async function download(version) {
    const arch = architecture();
    const name = executable();
    const rsp = await fetch(`https://dl.k8s.io/release/${version}/bin/${Deno.build.os}/${arch}/kubectl`);
    await Deno.writeFile(name, rsp.body);
}

const version = await latest();
if (await already(version)) {
    console.log("already installed:", version)
} else {
    console.log("downloading:", version);
    await download(version);
    console.log("done.");
}