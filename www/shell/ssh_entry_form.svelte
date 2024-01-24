<script>
    import { createEventDispatcher } from "svelte";
    import { equal } from "../utility.js";

    let ssh = {
        desc: "",
        host: "",
        port: 22,
        user: "root",
        pass: "",
    }, old = {};

    export function update(e) {
        old = e;
        for (const [key, val] of Object.entries(e)) {
            if (ssh[key] !== "undefined") ssh[key] = val;
        }
    }
    let inputHost;
    export function focus() {
        inputHost.select();
        inputHost.focus();
    }

    const dispatch = createEventDispatcher();

    function onSubmit() {
        onSave()
        dispatch("submit");
    }

    function onSave(e) {
        if (ssh.desc == "") ssh.desc = `${ssh.user}@${ssh.host}:${ssh.port}`
        dispatch("save", ssh);
    }
</script>

<table class="table">
    <thead><tr><th>编辑</th></tr></thead>
    <tbody><tr><td>
        <form on:submit|preventDefault={onSubmit}>
            <div class="row mb-2">
                <lable class="col-md-3 col-form-label">Host:</lable>
                <div class="col-md-9">
                    <input bind:this={inputHost} type="text" name="host" class="form-control" bind:value={ssh.host} />
                </div>
            </div>
            <div class="row mb-2">
                <lable class="col-md-3 col-form-label">Port:</lable>
                <div class="col-md-9">
                    <input type="number" name="port" class="form-control" bind:value={ssh.port} />
                </div>
            </div>
            <div class="row mb-2">
                <lable class="col-md-3 col-form-label">User:</lable>
                <div class="col-md-9">
                    <input type="text" name="user" class="form-control" bind:value={ssh.user} />
                </div>
            </div>
            <div class="row mb-2">
                <lable class="col-md-3 col-form-label">Pass:</lable>
                <div class="col-md-9">
                    <input type="password" name="pass" class="form-control" bind:value={ssh.pass} />
                </div>
            </div>
            <div class="row mb-2">
                <lable class="col-md-3 col-form-label">Desc:</lable>
                <div class="col-md-9">
                    <input type="text" name="desc" class="form-control" bind:value={ssh.desc} />
                </div>
            </div>
            <div class="row mb-2">
                <div class="offset-md-3 col-md-9">
                    <button type="button" disabled={equal(ssh, old)} on:click={onSave} class="btn btn-secondary"><i class="bi bi-check2-square"></i> 保存</button>
                    <button type="submit" on:click={onSubmit} class="btn btn-primary"><i class="bi bi-terminal"></i> 连接</button>
                </div>
            </div>
        </form>
    </td></tr></tbody>
</table>