import os from "node:os";

export function normalizeArch(arch) {
    if (arch == "x64" || arch == "x86_64") {
        return "amd64";
    }
    return arch
}

export function normalizePlatform(plat) {
    if (plat == "win32") {
        return "windows";
    }
    return plat;
}

export const arch = (function() {
    return normalizeArch(os.arch);
})();

export const platform = (function() {
    return normalizePlatform(os.platform());
})()

export function executable(name, platform, arch) {
    const ext = platform == "windows" ? ".exe" : "";
    return `bin/${name}_${platform}_${arch}${ext}`;
}
