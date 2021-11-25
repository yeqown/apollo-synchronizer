<template>
  <div id="app-wrapper">
    <a-layout style="height=100%">
      <a-layout-header>
        <h3>Asy tool</h3>
      </a-layout-header>
      <a-layout>
        <a-layout-sider theme="light" :collapsed="true">
          <a-menu mode="inline" style="height=100%">
            <a-menu-item key="upload"
              ><template #icon>
                <SyncOutlined two-tone-color="#1890ff" /> </template
              >Synchronize</a-menu-item
            >
            <a-menu-item key="setting"
              ><template #icon> <SettingTwoTone /> </template
              >Setting</a-menu-item
            >
          </a-menu>
        </a-layout-sider>
        <a-layout>
          <a-layout-content style="padding: 1em">
            <router-view></router-view>
          </a-layout-content>
          <a-layout-footer>copyright@yeqown 2021</a-layout-footer>
        </a-layout>
      </a-layout>
    </a-layout>
  </div>
</template>

<script>
import {
  Layout,
  LayoutHeader,
  LayoutSider,
  LayoutFooter,
  LayoutContent,
  Menu,
  MenuItem,
} from "ant-design-vue";
import { SettingTwoTone, SyncOutlined } from "@ant-design/icons-vue";
export default {
  name: "App",
  components: {
    ALayout: Layout,
    ALayoutHeader: LayoutHeader,
    ALayoutSider: LayoutSider,
    ALayoutContent: LayoutContent,
    ALayoutFooter: LayoutFooter,
    AMenu: Menu,
    AMenuItem: MenuItem,
    // icons
    SettingTwoTone,
    SyncOutlined,
  },
};
</script>

<style>
#body {
  height: 100%;
  width: 100%;
}

#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #ffffff;
  background-color: #e6f7ff;
  height: 100%;
  width: 100%;
}

#app-wrapper {
  height: 100%;
  width: 100%;
}
</style>
