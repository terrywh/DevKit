<script>
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    let timeout, inputFilter;

    function doSubmit(confirm) {
        clearTimeout(timeout);
        timeout = setTimeout(((confirm) => {
            return function() {
                dispatch("submit", {value: inputFilter.value, confirm: confirm});
            }
        })(confirm), 240);
    }

    function onKeydown(e) {
        doSubmit(false);
    }

    function onSubmit(e) {
        doSubmit(true);
    }

</script>

<form class="d-flex" on:submit|preventDefault={onSubmit}>
    <div class="flex-grow-1 me-2">
        <div class="input-group">
            <span class="input-group-text"><i class="bi bi-funnel"></i></span>
            <input type="text" bind:this={inputFilter} class="form-control" placeholder="搜索过滤" aria-label="filter" on:keydown={onKeydown} />
        </div>
    </div>
    <div>
        <button type="submit" class="btn btn-primary"><i class="bi bi-terminal"></i> 连接</button>
    </div>
</form>