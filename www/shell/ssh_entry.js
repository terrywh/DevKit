import { mount } from "svelte";
import SshEntry from "./ssh_entry.svelte"

const app = mount(SshEntry, {
    target: document.body,
});

export default app;