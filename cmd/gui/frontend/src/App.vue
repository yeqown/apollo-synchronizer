<template>
  <a-layout style="height: 100%">
    <div id="fix-window-placeholder" data-wails-drag>
      <h4>Apollo Synchronizer</h4>
    </div>
    <a-layout-header
      style="
        padding: 0;
        width: 100%;
        height: 48px;
        background: #ffffff;
        float: left;
      "
    >
      <div class="logo-container">
        <div id="logo">
          <img alt="logo" src="logo.png" style="width: 100%" />
        </div>
      </div>
      <a-menu
        mode="horizontal"
        theme="light"
        :default-selected-keys="['dashboard']"
        :style="{ lineHeight: '48px' }"
        @select="handleMenuSelect"
      >
        <a-menu-item key="dashboard" style="line-height: 48px">
          <template #icon> <DashboardOutlined /> </template>Welcome</a-menu-item
        >
        <a-menu-item key="synchronize" style="line-height: 48px"
          ><template #icon> <SyncOutlined two-tone-color="#1890ff" /> </template
          >Synchronize</a-menu-item
        >
        <a-menu-item key="setting" style="line-height: 48px"
          ><template #icon> <SettingOutlined /> </template>Setting</a-menu-item
        >
      </a-menu>
    </a-layout-header>
    <a-layout>
      <a-layout-content style="padding: 1em">
        <router-view></router-view>
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script>
import {
  Layout,
  LayoutHeader,
  // LayoutSider,
  // LayoutFooter,
  LayoutContent,
  Menu,
  MenuItem,
  // Affix,
} from "ant-design-vue";
import {
  SettingOutlined,
  SyncOutlined,
  DashboardOutlined,
} from "@ant-design/icons-vue";
export default {
  name: "App",
  components: {
    ALayout: Layout,
    ALayoutHeader: LayoutHeader,
    // ALayoutSider: LayoutSider,
    ALayoutContent: LayoutContent,
    // ALayoutFooter: LayoutFooter,
    AMenu: Menu,
    AMenuItem: MenuItem,
    // AAffix: Affix,
    // icons
    SettingOutlined,
    SyncOutlined,
    DashboardOutlined,
  },
  data() {
    return {
      settings: [
        {
          name: "setting1",
          env: "dev", // should be unique
          data: {},
        },
      ],
    };
  },
  methods: {
    handleMenuSelect({ item, selectedKeys }) {
      // console.log(item, key, selectedKeys);
      if (selectedKeys.length !== 0 && item && !item.disabled) {
        this.$router.push(`/${selectedKeys[0]}`);
      }
    },
  },
};
</script>

<style>
#body {
  /* width: 1280px;
  height: 740px; */
  width: 800px;
  height: 600px;
}

#app {
  font-family: "PT Mono", Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  height: 100%;
  width: 100%;
}

.logo-container {
  float: left;
  width: 120px;
  height: 48px;
  display: flex;
  flex-direction: column;
}

#fix-window-placeholder {
  height: 36px;
  width: 100%;
  background: #1890ff;
}

#fix-window-placeholder > h4 {
  font-weight: bold;
  line-height: 36px;
  margin: 0px;
  color: azure;
}

.logo-container > #logo {
  width: 100%;
  height: 48px;
  line-height: 48px;
  background: #fbfbfe;
}
</style>
