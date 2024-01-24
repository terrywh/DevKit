import { createRoot } from "svelte";
import TkeEntry from "./tke_entry.svelte"

const app = createRoot(TkeEntry, {
    target: document.body,
});

export default app;