<script>
    import { tick } from "svelte";
    import { route } from "../store.js";
    import { tkeEntry } from "./store.js";
    import SshEntryFilter from "./ssh_entry_filter.svelte";
    import TkeEntryForm from "./tke_entry_form.svelte";
    import TkeEntryJump from "./tke_entry_jump.svelte";
    import TkeEntryList from "./tke_entry_list.svelte";

    let connectForm, entryList, entryForm;
    let connectWindowTarget = $state("");

    $effect(() => {
        const index = $route.get("entry", 0);
        connectWindowTarget = `shell-tke-${index}`;
    });

    let entryFilter = $state("");

    function onJumpSubmit(e) {
        console.log("jump submit:", e.detail);
    }

    function onListSelect(e) {
        // entryForm.focus();
    }

    function onListSubmit(e) {
        doConnect();
    }

    function onFormSubmit(e) {
        console.log("form submit: ", e.detail);
        const index = $tkeEntry.append($route.get("entry", 0), e.detail);
        $route.put("entry", index);
    }

    async function onFilterSubmit(e) {
        entryList.$set({filter: e.detail.value});
        await tick();
        if (e.detail.confirm) doConnect();
    }

    function doConnect() {
        connectForm.submit();
    }

</script>


<div class="container mt-2">
    <form bind:this={connectForm} target={connectWindowTarget} action="/shell/shell.html">
        <input type="hidden" name="entry" value={$route.get("entry")} />
        <input type="hidden" name="type" value="tke" />
    </form>
    <div class="row">
        <div class="col-12">
            <TkeEntryJump on:submit={onJumpSubmit}></TkeEntryJump>
        </div>
    </div>
    <div class="row mb-2">
        <div class="col-12">
            <TkeEntryForm bind:this={entryForm} on:submit={onFormSubmit}></TkeEntryForm>
        </div>
    </div>
    <div class="row mb-2">
        <div class="col-12">
            <SshEntryFilter  bind:filter={entryFilter} onsubmit={onListSubmit}></SshEntryFilter>
        </div>
    </div>
    <div class="row mb-2">
        <div class="col-12">
            <TkeEntryList bind:this={entryList} bind:filter={entryFilter} on:select={onListSelect} on:submit={onListSubmit}></TkeEntryList>
        </div>
    </div>
    
</div>