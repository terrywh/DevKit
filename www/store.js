// import { subscribe } from "svelte/store";

export function createListStorage(key, def, sort) {
    const _proxy = {
        store: [], // 实际数据
        $subscriber: new Set(),
        subscribe(update, invalidate) {
            const subscriber = [update, invalidate];
            _proxy.$subscriber.add(subscriber);
            update(_proxy);
            return {
                unsubscribe() {
                    _proxy.$subscriber.delete(subscriber);
                },
            }
        },
        load() {
            this.store = JSON.parse(localStorage.getItem(key));
            if (!this.store) this.store = def ? [def] : [];
            this._update();
        },
        save() {
            // if (sort) this.store = this.store.sort(sort);
            localStorage.setItem(key, JSON.stringify(this.store));
            this._update();
        },
        _update() {
            for (const [update, invalidate] of this.$subscriber) update(this); // 通知所有订阅者
        },
        append(index, value) {
            value = Object.assign({}, value);
            index = parseInt(index) || 0;
            if (index == 0) {
                index = this.store.length;
                this.store.push(value);
            } else {
                this.store[index] = value;
            }
            this.save();
            return index;
        },
        fetch(index) {
            return this.store ? this.store[index] : undefined;
        },
        remove(index) {
            this.store.splice(index, 1);
            this.save();
            if (index < this.store.length) return index;
            else return this.store.length - 1;
        }
    }
    _proxy.load();
    return _proxy;
}

const _route = {
    store: new URLSearchParams(window.location.hash.substring(1)),
    $subscriber: new Set(),
    load() {
        this.store = new URLSearchParams(window.location.hash.substring(1));
        this._update();
    },
    _onhashchange() {
        if (!_route._putting) _route.load();
        _route._putting = false;
    },
    _update() {
        for (const [update, invalidate] of this.$subscriber) update(this); // 通知所有订阅者
    },
    subscribe(update, invalidate) {
        const subscriber = [update, invalidate];
        _route.$subscriber.add(subscriber);
        update(_route);
        return {
            unsubscribe() {
                _route.$subscriber.delete(subscriber);
            },
        }
    },
    get(key, def) {
        if (this.store.has(key)) {
            return this.store.get(key);
        } else {
            return def;
        }
    },
    put(key, val) {
        if (this.store.get(key) != val) {
            this.store.set(key, val);
            history.replaceState(null, null, "#" + this.store.toString()); // 不要触发 hashchange 事件（循环更新）
            this._update();
        }
    },
    getAll(key) {
        return this.store.getAll(key)
    },
    putAll(key, vals) {
        this.store.delete(key);
        for (const val of vals) {
            this.store.append(key, val);
        }
        history.replaceState(null, null, "#" + this.store.toString()); // 不要触发 hashchange 事件（循环更新）
        this._update();
    },
    onLink(e) {
        e.preventDefault();
        const r = new URLSearchParams(e.target.getAttribute("href").substring(1));
        for (const [key, val] of r.entries()) {
            _route.put(key, val)
        }
    }
}

window.addEventListener("hashchange", _route._onhashchange)
_route.load();
export const route = {
    subscribe: _route.subscribe,
};
