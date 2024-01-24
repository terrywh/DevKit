import { createRoot } from "svelte";
import SshEntry from "./ssh_entry.svelte"

const app = createRoot(SshEntry, {
    target: document.body,
});

export default app;