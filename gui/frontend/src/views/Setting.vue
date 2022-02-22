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
          <a-list-item key="item.title">
            <template #extra>
              <img
                width="272"
                alt="logo"
                src="https://gw.alipayobjects.com/zos/rmsportal/mqaQswcyDLcXyDKnZfES.png"
              />
            </template>
            <a-list-item-meta :description="item.description">
              <template #title>
                <a :href="item.href">{{ item.title }}</a>
              </template>
              <template #avatar><a-avatar :src="item.avatar" /></template>
            </a-list-item-meta>
            <!-- data -->
            {{ item.content }}
          </a-list-item>
        </template>
      </a-list>
    </div>
  </div>
</template>

<script>
import {
  PageHeader,
  Button,
  List,
  ListItem,
  ListItemMeta,
  Empty,
  Avatar,
} from "ant-design-vue";
export default {
  name: "Setting",
  components: {
    APageHeader: PageHeader,
    AButton: Button,
    AList: List,
    AListItem: ListItem,
    AListItemMeta: ListItemMeta,
    AEmpty: Empty,
    AAvatar: Avatar,
  },
  data() {
    return {
      settings: [
        {
          title: "setting1",
          description: "description1",
          cluster: "dev",
          avatar:
            "https://zos.alipayobjects.com/rmsportal/ODTLcjxAfvqbxHnVXCYX.png",
          data: {},
        },
        {
          title: "setting2",
          description: "description2",
          cluster: "prod",
          avatar:
            "https://zos.alipayobjects.com/rmsportal/ODTLcjxAfvqbxHnVXCYX.png",
          data: {},
        },
      ],
    };
  },
};
</script>