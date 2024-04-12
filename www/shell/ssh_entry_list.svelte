<script>
    import { createEventDispatcher, tick } from "svelte";
    import { route } from "../store.js";
    import { sshEntry } from "./store.js";

    const dispatch = createEventDispatcher();
    let { filter = $bindable("") } = $props();

    function onSelect(index, e) {
        if (!$sshEntry.store[index]) return;
        $route.put("entry", index);
        dispatch("select", {index, click: !!e});
    }

    function onDelete(index) {
        if (index < 1) return;
        $route.put("entry", $sshEntry.remove(index));
    }

    function filterApply(e, f) {
        if (!e || !f) return true;
        return e.host && e.host.indexOf(f) > -1 || e.desc && e.desc.indexOf(f) > -1;
    }

    function filterSelect(f) {
        if (!f) return;
        for(let i=0;i<$sshEntry.store.length;++i) {
            if (filterApply($sshEntry.store[i], filter)) {
                onSelect(i);
                break; // 找到第一个选中项即可
            }
        }
    }

    $effect(() => {
        filterSelect(filter);
    });

    function onSubmit(index) {
        dispatch("submit", {index});
    }
</script>

<table class="table">
    <thead>
        <tr>
            <th>#</th>
            <th>备注</th>
            <th>HOST</th>
            <th>PORT</th>
            <th>USER</th>
            <th style="width:10rem;">操作</th>
        </tr>
    </thead>
    <tbody>
    {#each $sshEntry.store as e, i}
        {#if filterApply(e, filter)}
        <tr class:table-primary={$route.get("entry", 0) == i} on:click={() => onSelect(i, true)}>
            <td>{i}</td>
            <td>{e.desc || "-"}</td>
            <td>{e.host || "-"}</td>
            <td>{e.port || "-"}</td>
            <td>{e.user || "-"}</td>
            <td>
                <div class="btn-group">
                    <button class="btn btn-secondary" disabled={i==0} title="跳板" on:click|preventDefault|stopPropagation={() => onSubmit(i)}><i class="bi bi-box-arrow-right"></i></button>
                    <button class="btn btn-secondary" disabled={i==0} title="删除" on:click|preventDefault|stopPropagation={() => onDelete(i)}><i class="bi bi-trash"></i></button>
                </div>
            </td>
        </tr>
        {/if}
    {/each}
    </tbody>
</table>
