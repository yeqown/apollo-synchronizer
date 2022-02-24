function loadSetting() {
    if (window.go && window.go.backend && window.go.backend.App && window.go.backend.App.LoadSetting) {
        return window.go.backend.App.LoadSetting();
    }

    // return Promise.resolve([]);
    // mock
    return Promise.resolve([
        {
            title: "setting1",
            account: "apollo",
            clusters: ["default", "swimming1"],
            envs: ["DEV"],
            portalAddr: "http://localhost:8080",
            secret: "ebba7e6efa4bb04479eb38464c0e7afc65",
            fs: "/Users/jia/.asy/setting1-DEV-$portalHash6",
        },
        {
            title: "setting2",
            account: "apollo",
            clusters: ["default", "preprod"],
            envs: ["DEV"],
            portalAddr: "http://localhost:8080",
            secret: "ebba7e6efa4bb04479eb38464c0e7afc65",
            fs: "/Users/jia/.asy/setting2-DEV-$portalHash6",
        },
    ]);
}

function saveSetting(settings) {
    if (window.go && window.go.backend && window.go.backend.App && window.go.backend.App.SaveSetting) {
        window.go.backend.App.SaveSetting(settings);
        return
    }

    return Promise.reject("No go.backend.App.SaveSetting loaded");
}

export { loadSetting, saveSetting };