import os from "node:os";

export const arch = (function() {
    const arch = os.arch();
    if (arch == "x64") {
        return "amd64";
    }
    return arch;
})();

export const platform = (function() {
    const plat = os.platform();
    if (plat == "win32") {
        return "windows";
    }
    return plat;
})()

export const ext = (function() {
    return platform == "windows" ? ".exe" : "";
}) ()

export const kubectl = (function() {
    return `bin/kubectl_${platform}_${arch}${ext}`;
})()

export const devkit = (function() {
    return `bin/devkit_${platform}_${arch}${ext}`;
}) ()

export const devkitServer = (function() {

}) ()