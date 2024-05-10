import { $ } from "bun";
import { parseArgs } from "node:util";
import { devkit, kubectl } from "./core";

const {options, target} = (function () {
    const r = parseArgs({
        args: Bun.argv.slice(2),
        options: {
            help: {
                type: "boolean",
                short: "h"
            },
        },
        strict: false,
        allowPositionals: true,
    });
    return {options: r.values, target: r.positionals};
})()


await (async function() {
    if (options.help) {
        console.log(Bun.argv.slice(0, 2).join(" "), "<target>")
        return
    }
    if (target[0].startsWith("cls-")) {
        $`./${kubectl} --kubeconfig=~/.kube/${target}.kubeconfig cp ${devkit} ${target[1]}:/root/devkit`
        $`./${kubectl} --kubeconfig=~/.kube/${target}.kubeconfig cp ${devkitServer} ${target[1]}:/root/devkit-server`
    }
}) ()