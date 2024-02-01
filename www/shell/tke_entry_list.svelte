<script>
    import { createEventDispatcher, tick } from "svelte";
    import { route } from "../store.js";
    import { tkeEntry, parseTkeInit } from "./store.js";

    const dispatch = createEventDispatcher();

    let { filter = "" } = $props();

    function filterApply(e, f) {
        if (!e || !f) return true;
        return e.host && e.host.indexOf(f) > -1 || e.desc && e.desc.indexOf(f) > -1;
    }

    function filterSelect(f) {
        if (!f) return;
        for(let i=0;i<$tkeEntry.store.length;++i) {
            if (filterApply($tkeEntry.store[i], filter)) {
                onSelect(i);
                break; // 找到第一个选中项即可
            }
        }
    }

    $effect(() => {
        filterSelect(filter);
    });

    async function onSelect(index) {
        $route.put("entry", index);
        await tick();
        dispatch("select", {index});
    }

    async function onSubmit(index) {
        $route.put("entry", index);
        await tick();
        dispatch("submit", {index})
    }

    function onDelete(index) {
        $route.put("entry", $tkeEntry.remove(index));
    }

</script>
<table class="table">
    <thead><tr>
        <th>#</th>
        <th>DESC</th>
        <th>CLUSTER</th>
        <th>NAMESPACE</th>
        <th>POD</th>
        <th>CONTAINER</th>
        <th></th>
    </tr></thead>
    <tbody>
        {#each $tkeEntry.store as e, i}
        {#if filterApply(e, filter)}
        {@const entry = parseTkeInit(e.init)}
        <tr on:click={() => onSelect(i)} class:table-primary={i == $route.get("entry")}>
            <td>{i}</td>
            <td>{e.desc}</td>
            <td>{entry.cluster}</td>
            <td>{entry.namespace}</td>
            <td>{entry.pod}</td>
            <td>{entry.container}</td>
            <td>
                <div class="btn-group">
                    <button class="btn btn-secondary" disabled={i==0} title="连接" on:click|preventDefault|stopPropagation={() => onSubmit(i)}><i class="bi bi-terminal"></i></button>
                    <button class="btn btn-secondary" disabled={i==0} title="删除" on:click|preventDefault|stopPropagation={() => onDelete(i)}><i class="bi bi-trash"></i></button>
                </div>
            </td>
        </tr>
        {/if}
        {/each}
    </tbody>
</table>