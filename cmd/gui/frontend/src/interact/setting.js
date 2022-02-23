function loadSetting() {
    return window.go.main.App.LoadSetting();
}

function saveSetting(settings) {
    return window.go.main.App.SaveSetting(settings);
}

export default { loadSetting, saveSetting };