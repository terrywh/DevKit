import { createRoot } from "svelte";
import Entry from "./k8s_entry.svelte"

const app = createRoot(Entry, {
    target: document.body,
});

export default app;