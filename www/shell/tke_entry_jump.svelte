<script>
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();
    let jump = localStorage.getItem("tke:store:jump") || "terryhaowu@csig.mnet2.com:36000";
    let temp = {
        addr: "",
        pwd1: "",
        pwd2: "",
    };
    function buildTemp() {
        // TODO 优化 URL 兼容健壮性
        const url = new URL("http://" + jump);
        temp.addr = `${url.username}@${url.host}`;
        temp.pwd1 = url.password.substring(0,6);
        temp.pwd2 = url.password.substring(6);
    }
    buildTemp();

    function onSubmit(e) {
        console.log("on submit");
        const tmp = buildJump(temp);
        const confirm = jump == tmp;
        dispatch("submit", {server: jump, confirm});
        localStorage.setItem("tke:store:jump", jump = tmp);
    }
    
    function buildJump(temp) {
        const tmp = temp.addr.split("@");
        return `${tmp[0]}:${temp.pwd1}${temp.pwd2}@${tmp[1]}`;
    }
    
</script>

<form class="d-flex" on:submit|preventDefault={onSubmit} >
    <div class="flex-grow-1">
        <div class="me-2 mb-2">
            <div class="input-group">
                <span class="input-group-text">地址：</span>
                <input class="form-control" type="text" bind:value={temp.addr} />
            </div>
        </div>
    </div>
    <div>
        <div class="float-start me-2 mb-2">
            <div class="input-group">
                <span class="input-group-text">密钥：</span>
                <input class="form-control" type="password" style="width: 10rem;" bind:value={temp.pwd1} />
            </div>
        </div>
        <div class="float-start me-2 mb-2">
            <div class="input-group">
                <span class="input-group-text">动态：</span>
                <input class="form-control" type="password" style="width: 10rem;" bind:value={temp.pwd2} />
            </div>
        </div>
    </div>
    <div>
        <button class="btn btn-secondary" disabled={buildJump(temp) == jump} type="submit"><i class="bi bi-check2-square"></i> 保存</button>
    </div>
</form>