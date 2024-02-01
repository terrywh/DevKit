<script>
    import { route } from "../store.js";
    import { k8sEntry } from "./store.js";

    function onSelect(index) {
        $route.put("entry", index);
    }
    function onDelete(index) {
        $k8sEntry.remove(index);
        if ($k8sEntry.store.length <= index) {
            index = $k8sEntry.store.length - 1;
            $route.put("entry", index);
        }
    }
</script>

<div class="list-group">
    {#each $k8sEntry.store as e, i}
    <a href="#index={i}" class="list-group-item list-group-item-action" class:list-group-item-primary={$route.get("entry", 0)==i}
        on:click|preventDefault={() => onSelect(i)}>
        {(e.desc || e.namespace)}
        <small class="text-secondary">{e.cluster_id && " @ " + e.cluster_id}</small>
        {#if i > 0}
        <div class="float-end btn-group mt-1">
            <button class="btn btn-light" title="删除" on:click|preventDefault|stopPropagation={() => onDelete(i)}><i class="bi bi-trash"></i></button>
        </div>
        {/if}
    </a>
    {/each}
</div>     
