<script>
    export let key = "";
    let configuring = 0, copying = 0;

    async function onConfig(e) {
        const form = new URLSearchParams(location.search);
        configuring = 1;    
        fetch("/bash/config?key=" + key).then(() => {
            configuring = 2
        }, () => {
            configuring = 3;
        });
    }

    let copyTimeout;
    async function onCopy(e) {
        await navigator.clipboard.writeText(document.title);
        copying = 1;
        clearTimeout(copyTimeout);
        copyTimeout = setTimeout(() => {
            copying = 0;
        }, 480);
    }

    function onEnter(e) {
        e.target.classList.remove("opacity-25")
    }
    function onLeave(e) {
        e.target.classList.add("opacity-25")
    }

</script>

<div class="opacity-25" role="toolbar" tabindex="-1" on:mouseenter={onEnter} on:mouseleave={onLeave} >
    <div class="btn-group float-start">
        <button type="button" class="btn btn-secondary" on:click={onConfig}>
            {#if configuring == 1}
            <div class="spinner-border" style="height: 1rem; width: 1rem;" role="status"></div>
            {:else if configuring == 2}
            <i class="bi bi-cloud-upload-fill"></i>
            {:else}
            <i class="bi bi-cloud-upload"></i>
            {/if}
        </button>
        <button type="button" class="btn btn-secondary" on:click={onCopy}>
            {#if copying == 1}
            <i class="bi bi-clipboard-check"></i>
            {:else}
            <i class="bi bi-clipboard"></i>
            {/if}
        </button>
    </div>
</div>