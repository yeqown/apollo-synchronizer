<template>
  <div>
    <a-page-header
      style="
        border: 1px solid rgb(235, 237, 240);
        margin-bottom: 1em;
        background-color: #ffffff;
      "
      title="Setting"
      sub-title="manages apollo cluster information and more."
    >
      <template #extra>
        <a-button key="add"> Add </a-button>
        <a-button key="add" :disabled="!modified"> Save </a-button>
      </template>
    </a-page-header>

    <!-- setting render -->
    <div style="height: 100%; text-align: left; background: #ffffff">
      <!-- empty -->
      <a-empty
        v-if="!settings || settings.length === 0"
        image="https://gw.alipayobjects.com/mdn/miniapp_social/afts/img/A*pevERLJC9v0AAAAAAAAAAABjAQAAAQ/original"
        :image-style="{
          height: '60px',
        }"
      >
        <template #description>
          <span> How about adding some settings? </span>
        </template>
        <a-button type="primary">Add Now</a-button>
      </a-empty>

      <!-- list -->
      <a-list item-layout="vertical" size="small" :data-source="settings">
        <template #renderItem="{ item }">
          <a-list-item :key="item.title">
            <!-- config data -->
            <a-descriptions
              :title="item.title"
              bordered
              size="small"
              :column="{ xxl: 4, xl: 3, lg: 3, md: 3, sm: 2, xs: 1 }"
            >
              <a-descriptions-item label="Portal">{{
                item.portalAddr
              }}</a-descriptions-item>
              <a-descriptions-item label="Clusters">{{
                item.clusters.join(",")
              }}</a-descriptions-item>
              <a-descriptions-item label="Account">
                {{ item.account }}
              </a-descriptions-item>
              <a-descriptions-item label="Secret">{{
                item.secret
              }}</a-descriptions-item>
              <a-descriptions-item label="Local Directory">{{
                item.fs
              }}</a-descriptions-item>
            </a-descriptions>
          </a-list-item>
        </template>
      </a-list>
    </div>
  </div>
</template>

<script>
import { loadSetting, saveSetting } from "../interact/index";
import {
  PageHeader,
  Button,
  List,
  ListItem,
  Empty,
  // Avatar,
  Descriptions,
  DescriptionsItem,
} from "ant-design-vue";
export default {
  name: "Setting",
  components: {
    APageHeader: PageHeader,
    AButton: Button,
    AList: List,
    AListItem: ListItem,
    AEmpty: Empty,
    // AAvatar: Avatar,
    ADescriptions: Descriptions,
    ADescriptionsItem: DescriptionsItem,
  },
  data() {
    return {
      modified: false,
      settings: [
        {
          title: "setting1",
          account: "apollo",
          clusters: ["default", "swimming1"],
          env: "DEV",
          portalAddr: "http://localhost:8080",
          secret: "ebba7e6efa4bb04479eb38464c0e7afc65",
          fs: "/Users/jia/.asy/setting1-DEV-$portalHash6",
        },
        {
          title: "setting2",
          account: "apollo",
          clusters: ["default", "preprod"],
          env: "DEV",
          portalAddr: "http://localhost:8080",
          secret: "ebba7e6efa4bb04479eb38464c0e7afc65",
          fs: "/Users/jia/.asy/setting2-DEV-$portalHash6",
        },
      ],
    };
  },
  mounted() {
    console.log("mounted");
    loadSetting().then((settings) => {
      console.log("loadSetting result", settings);
      this.settings = settings;
    });
  },
  methods: {
    enableModified() {
      this.modified = true;
    },
    save() {
      saveSetting(this.settings).then((result) => {
        console.log(this.settings);
        console.log("SaveSettings called: ", result);
      });
      this.modified = false;
    },
  },
};
</script>