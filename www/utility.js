
export async function hash() {
    const e = new TextEncoder();
    let o = new Uint8Array(0), p;
    
    for (let v of arguments) {
        if (typeof v !== "string") v = JSON.stringify(v);
        v = e.encode(v);
        p = new Uint8Array(o.byteLength + v.byteLength);
        p.set(o, 0);
        p.set(v, o.byteLength);
        o = p;
    }

    const buffer = await crypto.subtle.digest("SHA-1", o);
    const data = Array.from(new Uint8Array(buffer));
    return data.map((byte) => byte.toString(16).padStart(2, "0")).join("");
}

export function equal(o1, o2) {
    const keys1 = Object.keys(o1 || {});
    const keys2 = Object.keys(o2 || {});

    if (keys1.length !== keys2.length) {
        return false;
    }

    for (let key of keys1) {
        if (o1[key] !== o2[key]) {
            return false;
        }
    }

    return true;
}

export function deepEqual(o1, o2) {
    const keys1 = Object.keys(o1);
    const keys2 = Object.keys(o2);

    if (keys1.length !== keys2.length) {
        return false;
    }

    for (const key of keys1) {
        const val1 = o1[key];
        const val2 = o2[key];
        const areObjects = isObject(val1) && isObject(val2);
        if (
            areObjects && !deepEqual(val1, val2) ||
            !areObjects && val1 !== val2
        ) {
            return false;
        }
    }

    return true;
}

export function isObject(object) {
    return object != null && typeof object === 'object';
}