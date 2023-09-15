<script>
    import { route } from "../store.js";
    import { hash } from "../utility.js";
    import { onMount } from "svelte";

    let cluster = {cluster_id: "", svc: [], pod: [], node: []};
    let load_;

    async function load() {
        clearTimeout(load_);
        load_ = setTimeout(async () => {
            const query = new URLSearchParams(window.location.search.substring(1));
            const rsp = await fetch("/cluster/describe?" + query.toString());
            const tmp = await rsp.json()
            for(const pod of tmp.pod) {
                pod.target = await hash(tmp.cluster_id, tmp.namespace, pod.name);
            }
            cluster = tmp;
        }, 100);
    }

    $: onLoad($route);

    function onLoad($route) {
        load()
    }

    onMount(() => {
        if (!$route.get("page")) {
            $route.put("page", "pod");
        }
    })
</script>


<div class="container">
    <div class="row mt-2">
        <div class="col-12">
            <ul class="nav nav-tabs">
                <li class="nav-item">
                  <a class="nav-link" class:active={$route.get("page") == "pod"} on:click={$route.onLink} href="#page=pod">Pod</a>
                </li>
                <li class="nav-item">
                  <a class="nav-link" class:active={$route.get("page") == "svc"} on:click={$route.onLink} href="#page=svc">Svc</a>
                </li>
                <li class="nav-item">
                  <a class="nav-link" class:active={$route.get("page") == "node"} on:click={$route.onLink} href="#page=node">Node</a>
                </li>
              </ul>
        </div>
    </div>
    <div class="row" class:d-none={$route.get("page") != "pod"} >
        <div class="col-12">
            <table class="table table-striped">
                <thead>
                    <tr>
                      <th scope="col">Name</th>
                      <th scope="col">IP</th>
                      <th scope="col">Node</th>
                      <th scope="col">Status</th>
                      <th scope="col">-</th>
                    </tr>
                </thead>
                <tbody>
                    {#each cluster.pod as pod}
                    <tr>
                        <td>{pod.name}</td>
                        <td>{pod.ip}</td>
                        <td>{pod.node}</td>
                        <td>{pod.status}</td>
                        <td>
                            <a class="btn btn-primary"
                                target={pod.target}
                                href="/shell/shell.html?cluster_id={cluster.cluster_id}&namespace={cluster.namespace}&pod={pod.name}"
                                role="button">Bash</a>
                        </td>
                    </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>
    <div class="row" class:d-none={$route.get("page") != "node"} >
        <div class="col-12">
            <table class="table table-striped">
                <thead>
                    <tr>
                      <th scope="col">IP</th>
                      <th scope="col">Status</th>
                    </tr>
                </thead>
                <tbody>
                    {#each cluster.node as node}
                    <tr>
                        <td>{node.ip}</td>
                        <td>{node.status}</td>
                    </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>
    <div class="row" class:d-none={$route.get("page") != "svc"} >
        <div class="col-12">
            <table class="table table-striped">
                <thead>
                    <tr>
                      <th scope="col">Name</th>
                      <th scope="col">Type</th>
                      <th scope="col">ClusterIP</th>
                      <th scope="col">Port</th>
                    </tr>
                </thead>
                <tbody>
                    {#each cluster.svc as svc}
                    <tr>
                        <td>{svc.name}</td>
                        <td>{svc.type}</td>
                        <td>{svc.cluster_ip}</td>
                        <td>{svc.port}</td>
                    </tr>
                    {/each}
                </tbody>
            </table>
        </div>
    </div>
</div>