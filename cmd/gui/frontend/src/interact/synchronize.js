const modeMapping = {
    "upload": 1,
    "download": 2,
}

const containsKey = (optional, key) => {
    if (optional === undefined || !(optional instanceof Array)) {
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

const EVENT_RENDER_DIFF = "event.render.diff";
const EVENT_RENDER_RESULT = "event.render.result";
const EVENT_INPUT_DECIDE = "event.input.decide";


const decideMapping = {
    "confirm": 1,
    "cancel": 2,
}

const inputDecide = (decide) => {
    if (window.runtime && window.runtime.EventsEmit) {
        window.runtime.EventsEmit(EVENT_INPUT_DECIDE, decide);
        return
    }

    // DO NOTHING, just for mock
    console.warn("window.runtime.EventEmit unavailable!");
    return
}

const bindEventOnce = (event, cb) => {
    if (window.runtime && window.runtime.EventsOnce) {
        window.runtime.EventsOnce(event, cb);
        return
    }

    // DO NOTHING, just for mock
    console.warn("window.runtime.EventsOnce unavailable!");
    return
}

export {
    modeMapping, containsKey, synchronize,
    bindEventOnce,
    EVENT_RENDER_DIFF, EVENT_RENDER_RESULT, EVENT_INPUT_DECIDE,
    decideMapping, inputDecide
}