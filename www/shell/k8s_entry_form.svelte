<script>
    import { route } from "../store.js";
    import { k8sEntry } from "./store.js"
    let value = {
        cluster_id: "",
        namespace: "",
        desc: "",
    }
    let elClusterID;

    $: onRoute($route, $k8sEntry);

    function onRoute($route, $k8sEntry) {
        const index = $route.get("entry", 0) || 0;
        const e = $k8sEntry.store[index];
        if (!!e) for(const [key, val] of Object.entries(e)) {
            if (typeof value[key] !== "undefined") value[key] = val;
        }
        if (index == 0) {
            value.desc = "";
        }
        if (elClusterID) elClusterID.focus();
    }

    function onSubmit() {
        if (!value.desc) value.desc = value.namespace + "@" + value.cluster_id;
        
        let index = $route.get("entry", 0) || 0;
        index = $k8sEntry.append(index, value);
        $route.put("entry", index);
    }
</script>

<form class="container" target={value.namespace + '@' + value.cluster_id} action="/shell/k8s_cluster.html" on:submit={onSubmit}>
    <div class="row mb-2">
        <lable class="col-sm-2 col-form-label">ClusterID:</lable>
        <div class="col-sm-10">
            <input type="text" bind:this={elClusterID} name="cluster_id" class="form-control" bind:value={value.cluster_id} />
        </div>
    </div>
    <div class="row mb-2">
        <lable class="col-sm-2 col-form-label">Namespace:</lable>
        <div class="col-sm-10">
            <input type="text" name="namespace" class="form-control" bind:value={value.namespace} />
        </div>
    </div>
    <div class="row mb-2">
        <lable class="col-sm-2 col-form-label">Desc:</lable>
        <div class="col-sm-10">
            <input type="text" name="desc" class="form-control" bind:value={value.desc} />
        </div>
    </div>
    <div class="row mb-2">
        <button type="submit" class="btn btn-primary"><i class="bi bi-google-play"></i> 保存并查看</button>
    </div>
</form>