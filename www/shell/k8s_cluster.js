import { mount } from "svelte";
import K8SCluster from "./k8s_cluster.svelte"

const app = mount(K8SCluster, {
    target: document.body,
});

const query = new URLSearchParams(window.location.search.substring(1));
document.title = "集群 - " + (query.get("desc") || query.get("cluster_id"));

export default app;