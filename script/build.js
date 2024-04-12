import { compile } from "svelte/compiler";
import { join, extname, dirname } from "node:path";
import { copyFile, mkdir, open, readFile, readdir, stat, writeFile } from "node:fs/promises";

async function createImportMap() {
    // {
    //     "imports": {
    //         "svelte": "/node_modules/svelte/src/runtime/index.js",
    //         "svelte/internal": "/node_modules/svelte/src/runtime/internal/index.js",
    //         "svelte/internal/disclose-version": "/node_modules/svelte/src/runtime/internal/disclose-version/index.js"
    //     }
    // }
    const desc = await Bun.file("./node_modules/svelte/package.json").json();
    const map = {imports: {}};
    for (const entry in desc.exports) {
        const value = desc.exports[entry];
        const file = value.browser ? (value.browser.default || value.browser) : value.default;
        if (entry.at(0) == '.' && file) {
            map.imports[`svelte${entry.substring(1)}`] = `/node_modules/svelte${file.substring(1)}`;
        }
    }
    map.imports["esm-env"] = "/node_modules/esm-env/prod-browser.js";
    // map.imports["xterm"] = "/node_modules/xterm/lib/xterm.js";
    // map.imports["xterm-addon-webgl"] = "/node_modules/xterm-addon-webgl/lib/xterm-addon-webgl.js";
    // map.imports["xterm-addon-fit"] = "/node_modules/xterm-addon-fit/lib/xterm-addon-fit.js";
    // map.imports["trzsz"] = "/node_modules/trzsz/lib/trzsz.mjs";
    return new TextEncoder().encode(`\n        <script type="importmap">\n${JSON.stringify(map, null, "   ")}\n        </script>\n`)
}

async function build(src, dst, map) {
    for (const name of await readdir(src, {recursive: true})) {
        const source = join(src, name);
        const target = join(dst, name);
        
        const sourceStat = await stat(source);
        if (sourceStat.isDirectory()) continue;

        await mkdir(dirname(target), {recursive: true});
        console.log(source, target);

        if (source.endsWith(".svelte") || source.endsWith(".svelte.js")) {
            const file = await readFile(source);
            const m = compile(new TextDecoder().decode(file), {
                dev: true,
                css: "injected",
            });
            
            await writeFile(target, m.js.code);
        } else if (extname(source) == ".html") {
            const file = await readFile(source);
            const x = file.indexOf("<head>");
            const dst = await open(target, "w");
            await dst.write(file.subarray(0, x + 6));
            await dst.write(map);
            await dst.write(file.subarray(x + 6));
            await dst.close();
        } else {
            await copyFile(source, target);
        }
    }
}

await build("./www", "./public", await createImportMap());
console.timeEnd("build");
