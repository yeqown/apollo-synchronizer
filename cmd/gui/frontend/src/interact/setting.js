function loadSetting() {
    if (window.go && window.go.backend && window.go.backend.App && window.go.backend.App.LoadSetting) {
        return window.go.backend.App.LoadSetting();
    }

    return Promise.resolve([]);
}

function saveSetting(settings) {
    if (window.go && window.go.backend && window.go.backend.App && window.go.backend.App.SaveSetting) {
        window.go.backend.App.SaveSetting(settings);
    }

    return Promise.reject("No go.backend.App.SaveSetting loaded");
}

export { loadSetting, saveSetting };