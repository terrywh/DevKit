import { createListStorage } from "../store.js";

export const sshEntry = {
    subscribe: createListStorage("ssh:store:entry", {"port": 22, "user": "root", "desc":"新增"}).subscribe,
};

export const k8sEntry = {
    subscribe: createListStorage("k8s:store:entry", {"desc": "新增"}).subscribe
};

export const tkeEntry = {
    subscribe: createListStorage("tke:store:entry", {"desc": "新增"}).subscribe
};

export function parseTkeInit(init) {
    const e = {
        cluster: "-",
        namespace: "-",
        pod: "-",
        container: "-",
    };
    if (!init) return e;
    init.split(" -").forEach((x) => {
        const y = x.split(" ");
        switch(y[0]) {
        case "cls":
            e.cluster = y[1].trim();
            break
        case "n":
            e.namespace = y[1].trim();
            break;
        case "p":
            e.pod = y[1].trim();
            break;
        case "c":
            e.container = y[1].trim();
            break;
        }
    });
    return e;
}
