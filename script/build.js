// import { compile } from "https://esm.sh/svelte@4.2.9/compiler";
// import { compile } from "npm:svelte/compiler";
import { compile } from "svelte/compiler";
import { join, extname } from "node:path";


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
        const file = value.browser ? (value.browser.default || value.browser) : value.default;
        if (entry.at(0) == '.' && file) {
            map.imports[`svelte${entry.substring(1)}`] = `/node_modules/svelte${file.substring(1)}`;
        }
    }
    map.imports["esm-env"] = "/node_modules/esm-env/prod-browser.js";
    map.imports["xterm"] = "/node_modules/xterm/lib/xterm.js";
    map.imports["xterm-addon-webgl"] = "/node_modules/xterm-addon-webgl/lib/xterm-addon-webgl.js";
    map.imports["xterm-addon-fit"] = "/node_modules/xterm-addon-fit/lib/xterm-addon-fit.js";
    // map.imports["trzsz"] = "/node_modules/trzsz/lib/trzsz.mjs";
    return map;
}

async function build(src, dst, map) {
    for await (const entry of Deno.readDir(src)) {
        const source = join(src, entry.name);
        const target = join(dst, entry.name);
        if (entry.isDirectory) {
            await Deno.mkdir(target, {recursive: true});
            await build(source, target, map);
        } else if (entry.name.endsWith(".svelte") || entry.name.endsWith(".svelte.js")) {
            const file = await Deno.readFile(source);
            const m = compile(new TextDecoder().decode(file), {
                dev: true,
                css: "injected",
            });
            await Deno.writeTextFile(target, m.js.code);
        } else if (".html" === extname(entry.name)) {
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

await build("./www", "./public",
    new TextEncoder().encode(`<!DOCTYPE html><script type="importmap">${JSON.stringify(await map())}</script>\n`));
console.log("done.");