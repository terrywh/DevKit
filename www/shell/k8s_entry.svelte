<script>
    import EntryList from "./k8s_entry_list.svelte";
    import EntryForm from "./k8s_entry_form.svelte";
</script>

<div class="container">
    <div class="row mt-2">
        <div class="col-lg-8">
            <EntryForm></EntryForm>
        </div>
        <div class="col-lg-4">
            <EntryList></EntryList>
        </div>
    </div>
</div>