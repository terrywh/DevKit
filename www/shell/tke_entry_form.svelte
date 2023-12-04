<script>
    import { createEventDispatcher } from "svelte";
    import { route } from "../store.js";
    import { tkeEntry } from "./store.js";
    import { equal } from "../utility.js";

    const dispatch = createEventDispatcher();
    let entry = {
        desc: "",
        init: "",
    };

    $: onUpdate($route, $tkeEntry);

    function onUpdate($route, $tkeEntry) {
        const index = $route.get("entry", 0);
        const e = $tkeEntry.fetch(index);
        if (e) {
            entry = Object.assign(entry, e);
            entry = entry;
        }
    }

    function onSubmit() {
        // const confirm = equal(entry, $tkeEntry.fetch(index));
        dispatch("submit", entry);
    }

    let inputInit;
    export function focus() {
        inputInit.select();
        inputInit.focus();
    }

</script>
<form on:submit|preventDefault={onSubmit} class="d-flex">
    <div class="flex-grow-1 me-2">
        <div class="input-group">
            <span class="input-group-text">指令：</span>
            <input bind:this={inputInit} class="form-control" type="text" bind:value={entry.init} />
        </div>
    </div>
    <div class="me-2">
        <div class="input-group">
            <span class="input-group-text">备注：</span>
            <input class="form-control" type="text" style="width: 10rem;" bind:value={entry.desc} />
        </div>
    </div>
    <div>
        <button class="btn btn-secondary" disabled={equal(entry, $tkeEntry.fetch($route.get("entry",0)))} type="submit"><i class="bi bi-check2-square"></i> 保存</button>
    </div>
</form>