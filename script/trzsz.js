#!deno
import { DOMParser, Element } from "https://deno.land/x/deno_dom@v0.1.43/deno-dom-wasm.ts";

async function latest() {
    const rsp = await fetch("https://github.com/trzsz/trzsz-go/releases/latest");
    const doc = new DOMParser().parseFromString(await rsp.text(), "text/html");
    const vr = /\d+\.\d+(\.(\d+)(\.\d+)?)?/;
    const match = doc.title.match(vr);
    return match ? match[0].toString() : "-";
}

async function download(version, os, arch) {
    const rsp = await fetch(`https://github.com/trzsz/trzsz-go/releases/download/v${version}/trzsz_${version}_${os}_${arch}.tar.gz`);
    await Deno.writeFile(`var/trzsz_${version}_${os}_${arch}.tar.gz`, rsp.body);
}

async function already(version) {
    for await (const entry of Deno.readDir("var")) {
        if (entry.name.startsWith("trzsz_") && entry.name.indexOf(version) > 0) {
            return true;
        }
    }
    return false;
}

const version = await latest();
if (await already(version)) {
    console.log(`already installed: v${version}`)
} else {
    console.log(`downloading: v${version}`);
    await download(version, "linux", "x86_64");
    await download(version, "linux", "aarch64");
    console.log("done.");
}