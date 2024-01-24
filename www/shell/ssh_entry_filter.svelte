<script>
    let timeout, eFilter;
    
    let { filter = "" } = $props();

    function doSubmit() {
        clearTimeout(timeout);
        timeout = setTimeout(() => {
            filter = eFilter.value;
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
        <input type="text" bind:this={eFilter} class="form-control" placeholder="搜索过滤" aria-label="filter" on:keydown={onKeydown} />
    </div>
</form>