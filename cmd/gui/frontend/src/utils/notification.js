import { notification } from 'ant-design-vue';

const notificationWithIcon = (type, message, description) => {
    notification[type]({
        message: message,
        description: description,
        placement: 'bottomRight',
        duration: 2,
    });
}

const notificationError = (description) => {
    notificationWithIcon('error', 'Error', description);
}

const notificationSuccess = (description) => {
    notificationWithIcon('success', 'Success', description);
}

const notificationInfo = (description) => {
    notificationWithIcon('info', 'Info', description);
}

const notificationWarning = (description) => {
    notificationWithIcon('warning', 'Warning', description);
}

export { notificationWithIcon, notificationError, notificationSuccess, notificationInfo, notificationWarning };