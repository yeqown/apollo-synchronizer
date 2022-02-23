import { loadSetting, saveSetting } from './setting';

const loadStatistics = () => {
    if (window.go && window.go.backend && window.go.backend.App && window.go.backend.App.LoadStatistics) {
        return window.go.backend.App.LoadStatistics();
    }

    return Promise.resolve({
        lastOpenTs: new Date().getTime() / 1000,
        firstOpenTs: new Date().getTime() / 1000,
        openCount: 0,
        openTime: 0,

        uploadCount: 2,
        uploadFileCount: 0,
        uploadFileSize: 1001,
        uploadFailedCount: 0,

        downloadCount: 0,
        downloadFileCount: 0,
        downloadFileSize: 0,
        downloadFailedCount: 0,
    });
};

export { loadSetting, saveSetting, loadStatistics };