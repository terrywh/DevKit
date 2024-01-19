import { compile } from "svelte/compiler";
import { exists, mkdir, readdir, rm } from "node:fs/promises";
import { extname } from "node:path";

async function build(entries) {
    return Bun.build({
        root: "./www",
        entrypoints: entries,
        outdir: "./public",
        sourcemap: "external",
        splitting: true,
        naming: "[dir]/[name].[ext]",
        target: "browser",
        plugins:[
            {
                name: "svelte",
                async setup(build) {
                    build.onLoad({filter: /\.svelte$/}, async function(e) {
                        const origin = await Bun.file(e.path).text();
                        const code = compile(origin, {filename: e.path});
                        return {contents: code.js.code, loader: "js"};
                    });
                },
            }
        ]
    });
}

async function rcopy(src, dst, filter) {
    if (! await exists(`${dst}/`)) await mkdir(`${dst}/`);

    const files = await readdir(src, {withFileTypes: true});
    for (const file of files) {
        if (file.isDirectory()) {
            rcopy(`${src}/${file.name}`, `${dst}/${file.name}`, filter);
        } else if (filter(file)) {
            await Bun.write(`${dst}/${file.name}`, new Response(Bun.file(`${src}/${file.name}`)));
        }
    }
}

// await rm("./public", {recursive: true, force: true});

await rcopy("./www", "./public", (file) => {
    const ext = extname(file.name);
    return file.isFile() && [".css", ".html", ".png"].indexOf(ext) > -1;
});

const r = await build([
    // "./www/demo/example.js",
    "./www/shell/shell.js",
    "./www/shell/k8s_cluster.js",
    "./www/shell/ssh_entry.js",
    "./www/shell/k8s_entry.js",
    "./www/shell/tke_entry.js",
]);

if (r.success) {
    console.log("done.");
} else {
    console.error(r.logs);
}
