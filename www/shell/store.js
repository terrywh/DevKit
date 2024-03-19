import { createListStorage } from "../store.js";

const CreateEntryCompare = function(f1, f2, f3) {
    return function(entry1, entry2) {
        if (entry1.desc == "新增") return -1;
        if (entry2.desc == "新增") return  1;
        if (f1) if (entry1[f1] != entry2[f1]) return entry1[f1] < entry2[f1] ? -1 : 1;
        if (f2) if (entry1[f2] != entry2[f2]) return entry1[f2] < entry2[f2] ? -1 : 1;
        if (f3) if (entry1[f3] != entry2[f3]) return entry1[f3] < entry2[f3] ? -1 : 1;
        return 0;
    }
}

export const sshEntry = {
    subscribe: createListStorage("ssh:store:entry", {"port": 22, "user": "root", "desc":"新增"},
        CreateEntryCompare("desc", "host")).subscribe,
};

export const k8sEntry = {
    subscribe: createListStorage("k8s:store:entry", {"desc": "新增"},
        CreateEntryCompare("desc", "cluster_id", "namespace")).subscribe
};

export const tkeEntry = {
    subscribe: createListStorage("tke:store:entry", {"desc": "新增"},
        CreateEntryCompare("desc", "init")).subscribe
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
