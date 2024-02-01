<script>
    import { onMount, tick } from "svelte";
    import { route } from "../store.js";
    import { sshEntry } from "./store.js";

    import RoutePath from "./ssh_route_path.svelte"
    import EntryForm from "./ssh_entry_form.svelte";
    import EntryList from "./ssh_entry_list.svelte";
    import EntryFilter from "./ssh_entry_filter.svelte";
    
    let entryForm, entryList, entryFilter = "", routePath = [];

    function onListSubmit(e) {
        // connect.submit()
        $route.put("route", e.detail.index);
    }

    function onListSelect(e) {
        if (e.detail.click) entryForm.focus();
    }

    function onFormSubmit(e) {
        doTerminal()
    }

    function onFormSave(e) {
        const index = $sshEntry.append($route.get("entry", 0), e.detail);
        entryForm.update($sshEntry.fetch(index));
        $route.put("entry", index);
    }
    
    async function onFilterSubmit(e) {
        entryList.$set({filter: e.detail.value});
        await tick();
        if (e.detail.confirm) doTerminal();
    }
   
    function doTerminal() {
        connect.submit()
    }

    let connect, connectWindowTarget = "";
    $: {
        const index = $route.get("entry", 0);
        connectWindowTarget = `shell-ssh-${index}`;
        if (entryForm) entryForm.update($sshEntry.fetch(index));
    }

    onMount(() => {
        const index = $route.get("entry", 0);
        entryForm.update($sshEntry.fetch(index));
    })

</script>
<div class="container">
    <div class="row mb-3 mt-2">
        <div class="col-lg-12">
            <form bind:this={connect} target={connectWindowTarget} action="/shell/shell.html">
                <input type="hidden" name="route" value={$route.get("route")} />
                <input type="hidden" name="entry" value={$route.get("entry")} />
                <input type="hidden" name="type" value="ssh" />
            </form>
            <RoutePath bind:this={routePath}></RoutePath>
        </div>
    </div>
    <div class="row mb-3">
        <div class="col-12">
            <EntryFilter bind:filter={entryFilter} onsubmit={onFormSubmit}></EntryFilter>
        </div>
    </div>
    <div class="row mb-2">
        <div class="col-9">
            <EntryList bind:this={entryList} bind:filter={entryFilter} on:submit={onListSubmit} on:select={onListSelect}></EntryList>
        </div>
        <div class="col-3">
            <EntryForm bind:this={entryForm} on:submit={onFormSubmit} on:save={onFormSave}></EntryForm>
        </div>
    </div>
</div>