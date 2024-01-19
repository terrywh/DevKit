import { compile } from "npm:svelte/compiler";
import { join, extname } from "https://deno.land/std@0.212.0/path/mod.ts"


async function map() {
    // {
    //     "imports": {
    //         "svelte": "/node_modules/svelte/src/runtime/index.js",
    //         "svelte/internal": "/node_modules/svelte/src/runtime/internal/index.js",
    //         "svelte/internal/disclose-version": "/node_modules/svelte/src/runtime/internal/disclose-version/index.js"
    //     }
    // }
    const desc = JSON.parse(await Deno.readTextFile("./node_modules/svelte/package.json"));
    const map = {imports: {}};
    for (const entry in desc.exports) {
        const value = desc.exports[entry];
        const file = value.browser ? value.browser.default : value.default;
        if (entry.at(0) == '.' && file) {
            map.imports[`svelte${entry.substring(1)}`] = `/node_modules/svelte${file.substring(1)}`;
        }
    }
    return map;
}

async function run(src, dst, map) {
    for await (const entry of Deno.readDir(src)) {
        const source = join(src, entry.name);
        const target = join(dst, entry.name);
        if (entry.isDirectory) {
            await Deno.mkdir(target, {recursive: true});
            await run(source, target, map);
        } else if (".svelte" === extname(entry.name)) {
            const file = await Deno.readFile(source);
            const m = compile(new TextDecoder().decode(file), {dev: true});
            await Deno.writeTextFile(target, m.js.code);
        } else if (".html" === extname(entry.name)) {
            console.log(source, target);
            const fs = await Deno.open(source, {"read": true});
            const ft = await Deno.open(target, {"create": true, "write": true, "truncate": true});
            const w = await ft.writable.getWriter();
            await w.write(map);
            // await w.close();
            w.releaseLock()
            await fs.readable.pipeTo(ft.writable);
            // fs.close();
            // ft.close();
        } else {
            await Deno.copyFile(source, target);
        }
    }
}

run("./www", "./public",
    new TextEncoder().encode(`<!DOCTYPE html><script type="importmap">${JSON.stringify(await map())}</script>\n`));