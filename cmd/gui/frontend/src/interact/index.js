import { loadSetting, saveSetting } from './setting';

const loadStatistics = () => {
    if (window.go && window.go.backend && window.go.backend.App && window.go.backend.App.Statistics) {
        console.log('loadStatistics called');
        return window.go.backend.App.Statistics();
    }

    return Promise.resolve({
        lastOpenTs: new Date().getTime() / 1000,
        firstOpenTs: new Date().getTime() / 1000,
        openCount: 0,
        openTime: 0,

        uploadCount: 0,
        uploadFileCount: 0,
        uploadFileSize: 0,
        uploadFailedCount: 0,

        downloadCount: 0,
        downloadFileCount: 0,
        downloadFileSize: 0,
        downloadFailedCount: 0,
    });
};

export { loadSetting, saveSetting, loadStatistics };