const modeMapping = {
    "upload": 1,
    "download": 2,
}

const containsKey = (optional, key) => {
    if (optional || !(optional instanceof Array)) {
        return false
    }

    return optional.indexOf(key) !== -1
}

const synchronize = (scope) => {
    if (window.go && window.go.backend && window.go.backend.App && window.go.backend.App.Synchronize) {
        return window.go.backend.App.Synchronize(scope);
    }

    return Promise.reject("No go.backend.App.Synchronize loaded");
}

export { modeMapping, containsKey, synchronize }