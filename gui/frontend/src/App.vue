<template>
  <a-layout style="height: 100%">
    <a-affix :offset-top="top">
      <a-layout-header
        style="padding: 0; width: 100%; background: #ffffff; float: left"
      >
        <div class="logo-container">
          <div id="fix-window-placeholder"></div>
          <div id="logo" />
        </div>
        <a-menu
          mode="horizontal"
          theme="light"
          :default-selected-keys="['dashboard']"
          :style="{ lineHeight: '64px' }"
          @select="handleMenuSelect"
        >
          <a-menu-item key="dashboard">
            <template #icon> <DashboardTwoTone /> </template
            >Dashboard</a-menu-item
          >
          <a-menu-item key="synchronize"
            ><template #icon>
              <SyncOutlined two-tone-color="#1890ff" /> </template
            >Synchronize</a-menu-item
          >
          <a-menu-item key="setting"
            ><template #icon> <SettingTwoTone /> </template>Setting</a-menu-item
          >
        </a-menu>
      </a-layout-header>
    </a-affix>
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
  Affix,
} from "ant-design-vue";
import {
  SettingTwoTone,
  SyncOutlined,
  DashboardTwoTone,
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
    AAffix: Affix,
    // icons
    SettingTwoTone,
    SyncOutlined,
    DashboardTwoTone,
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

/* #components-layout-demo-top-side {
  width: 100%;
  height: 64px;
  padding: 0px;
  background: #ffffff;
} */

.logo-container {
  float: left;
  width: 120px;
  height: 64px;
  display: flex;
  flex-direction: column;
}

.logo-container > #fix-window-placeholder {
  height: 24px;
}

.logo-container > #logo {
  width: 100%;
  height: 40px;
  line-height: 40px;
  background: #e6f7ff;
}
</style>
