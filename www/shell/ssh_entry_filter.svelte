<script>
    let timeout;
    let { filter = $bindable(""), onsubmit = null } = $props();
    let filterInput = $state(filter);

    function doSubmit(submit) {
        clearTimeout(timeout);
        timeout = setTimeout(() => {
            filter = filterInput;
            if (submit && onsubmit) {
                onsubmit(filter);
            }
        }, 240);
    }

    function onKeydown(e) {
        doSubmit(false);
    }

    function onSubmit(e) {
        doSubmit(true);
    }

</script>

<form on:submit|preventDefault={onSubmit}>
    <div class="input-group">
        <span class="input-group-text"><i class="bi bi-funnel"></i></span>
        <input type="text" bind:value={filterInput} class="form-control" placeholder="搜索过滤" aria-label="filter" on:keydown={onKeydown} />
    </div>
</form>