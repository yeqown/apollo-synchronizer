import { createApp } from 'vue'
// import Antd from "ant-design-vue";
import App from './App.vue'
import "ant-design-vue/dist/antd.css"; // or 'ant-design-vue/dist/antd.less'
// import Wails from "@wailsapp/runtime";

const app = createApp(App)
app.config.productionTip = false

// app.use(Antd);
app.mount('#app')

// Wails.Init(() => {
//     createApp(App).mount('#app')
// })