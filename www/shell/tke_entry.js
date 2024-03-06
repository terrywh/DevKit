import { mount } from "svelte";
import TkeEntry from "./tke_entry.svelte"

const app = mount(TkeEntry, {
    target: document.body,
});

export default app;