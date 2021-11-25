import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import "ant-design-vue/dist/antd.css";
// import Antd from "ant-design-vue";
// import Wails from "@wailsapp/runtime";

const app = createApp(App)
app.config.productionTip = false

app.use(router)
app.mount('#app')