import { mount } from "svelte";
import Entry from "./k8s_entry.svelte"

const app = mount(Entry, {
    target: document.body,
});

export default app;