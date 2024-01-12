<script>
    import { createEventDispatcher } from "svelte";

    export let key = "";
    let configuring = 0, copying = 0;
    export let refreshing = false;

    const dispatch = createEventDispatcher();

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

    function onRefreshing(s) {
        console.log("onRefreshing: ", s);
        dispatch("refresh", {enable: s});
    }

    $: onRefreshing(refreshing);

</script>

<div class="opacity-25" role="toolbar" tabindex="-1" on:mouseenter={onEnter} on:mouseleave={onLeave} >
    <div class="btn-group float-start">
        <button type="button" class="btn btn-secondary" on:click={onConfig} title="安装文件工具">
            {#if configuring == 1}
            <div class="spinner-border" style="height: 1rem; width: 1rem;" role="status"></div>
            {:else if configuring == 2}
            <i class="bi bi-cloud-upload-fill"></i>
            {:else}
            <i class="bi bi-cloud-upload"></i>
            {/if}
        </button>
        <button type="button" class="btn btn-secondary" on:click={onCopy} title="复制标题">
            {#if copying == 1}
            <i class="bi bi-clipboard-check"></i>
            {:else}
            <i class="bi bi-clipboard"></i>
            {/if}
        </button>
        <input type="checkbox" class="btn-check" id="btn-refresh" bind:checked={refreshing} />
        <label class="btn btn-outline-primary" for="btn-refresh" title="自动保活">
            {#if refreshing}
            <div class="spinner-border" style="height: 1rem; width: 1rem;" role="status"></div>
            {:else}
            <i class="bi bi-arrow-clockwise"></i>
            {/if}
        </label>
    </div>
</div>