import { platform, arch, kubect, ext } from "./core.js";
import { $ } from "bun";

const latest = await (async function() {
    const rsp = await fetch("https://dl.k8s.io/release/stable.txt");
    return await rsp.text();
})()

const local = await (async function() {
    try {
        const json = await $`./${kubect} version --client=true -o json`.json()
        return json.clientVersion.gitVersion;
    } catch(e) {
        return "v0.0.0";
    }
})()

if (latest == local) {
    console.log(">", file)
    console.log("latest already installed:", latest)
} else {
    console.log(">", file)
    console.log("upgrading:", latest, "<=", local);
    await $`curl -L https://dl.k8s.io/release/${version}/bin/${platform}/${arch}/kubectl${ext} -o ${kubect}`;
    await $`chmod +x ${file}`;
    console.log("done.");
}
