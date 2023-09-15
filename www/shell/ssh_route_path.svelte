<script>
    import { route } from "../store.js";
    import { sshEntry } from "./store.js";
    let routes = [];
    let entry = {};
    $: {
        routes.splice(0, routes.length);
        for (const index of $route.getAll("route")) {
            routes.push(Object.assign({}, $sshEntry.fetch(index), {index}));
        }
        routes = routes;
    }

    $: {
        const index = $route.get("entry");
        if (index) entry = $sshEntry.fetch(index);
    }

    function onCancel(index) {
        const rs = [];
        for (const r of routes) {
            if (r.index == index) continue;
            rs.push(r.index);
        }
        $route.putAll("route", rs);
    }
    
</script>

<style>
    .route-start { padding: 0.8rem 0; }
    .route-joint { padding: 0.8rem 0.5rem; }
    .route-entry { padding: 0.8rem 1rem; background-color: #fafafa; }
</style>

<div>
    <div class="float-start route-start text-secondary">连接</div>
    {#each routes as r}
    <i class="float-start route-joint bi bi-chevron-right"></i>
    <div class="float-start route-entry round border">
        {`${r.desc} (${r.user}@${r.host})`}
        &nbsp;&nbsp; <a href="#remove" title="取消" on:click|preventDefault={() => onCancel(r.index)}><i class="bi bi-x-lg"></i></a>
    </div>
    {/each}
    <i class="float-start route-joint bi bi-chevron-right"></i>
    <div class="float-start route-entry round border">{`${entry.desc} (${entry.user}@${entry.host})`}</div>
</div>