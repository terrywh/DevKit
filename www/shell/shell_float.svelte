<script>
    let { key = "", refreshing = $bindable(false) } = $props();
    let configuring = 0, copying = 0;

    async function onConfig(e, arch) {
        const form = new URLSearchParams(location.search);
        configuring = 1;    
        fetch("/bash/config?key=" + key + "&arch=" + arch).then(() => {
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

    window.addEventListener("blur", function() {
        refreshing = true;
    });
    window.addEventListener("focus", function() {
        refreshing = false;
    });
</script>

<div class="opacity-25" role="toolbar" tabindex="-1" on:mouseenter={onEnter} on:mouseleave={onLeave} >
    <div class="btn-group float-start">
        <button type="button" class="btn btn-secondary" on:click={(e) => { onConfig(e, "x86_64") }} title="x86_64">
            {#if configuring == 1}
            <div class="spinner-border" style="height: 1rem; width: 1rem;" role="status"></div>
            <!-- {:else if configuring == 2}
            <i class="bi bi-cloud-upload"></i>
            X86 -->
            {:else}
            <i class="bi bi-cloud-upload"></i>
            X86
            {/if}
        </button>
        <button type="button" class="btn btn-secondary" on:click={(e) => { onConfig(e, "aarch64") }} title="aarch64">
            {#if configuring == 1}
            <div class="spinner-border" style="height: 1rem; width: 1rem;" role="status"></div>
            <!-- {:else if configuring == 2}
            <i class="bi bi-cloud-upload"></i>
            X86 -->
            {:else}
            <i class="bi bi-cloud-upload"></i>
            ARM
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