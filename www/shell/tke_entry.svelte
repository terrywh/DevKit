<script>
    import { tick } from "svelte";
    import { route } from "../store.js";
    import SshEntryFilter from "./ssh_entry_filter.svelte";
    import TkeEntryForm from "./tke_entry_form.svelte";
    import TkeEntryJump from "./tke_entry_jump.svelte";
    import TkeEntryList from "./tke_entry_list.svelte";
    import K8sEntryForm from "./k8s_entry_form.svelte";

    let connect, connectWindowTarget, entryList, entryForm;

    $: {
        const index = $route.get("entry", 0);
        connectWindowTarget = `shell-tke-${index}`;
    }

    function onJumpSubmit(e) {
        console.log("jump submit:", e.detail);
    }

    function onListSelect(e) {
        entryForm.focus();
    }

    function onFormSubmit(e) {
        console.log("form submit: ", e.detail);
    }

    async function onFilterSubmit(e) {
        entryList.$set({filter: e.detail.value});
        await tick();
        if (e.detail.confirm) doConnect();
    }

    function doConnect() {
        connect.submit();
    }

</script>


<div class="container mt-2">
    <form bind:this={connect} target={connectWindowTarget} action="/shell/shell.html">
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
            <SshEntryFilter on:submit={onFilterSubmit}></SshEntryFilter>
        </div>
    </div>
    <div class="row mb-2">
        <div class="col-12">
            <TkeEntryList bind:this={entryList} on:select={onListSelect}></TkeEntryList>
        </div>
    </div>
    
</div>