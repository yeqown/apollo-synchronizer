<template>
  <div>
    <a-page-header
      title="Synchronize"
      sub-title="upload or download configs between apollo and local fs."
      style="
        border: 1px solid rgb(235, 237, 240);
        margin-bottom: 1em;
        background-color: #ffffff;
      "
    ></a-page-header>
    <div class="card-container">
      <!-- upload -->
      <div class="card">
        <a-cloud-upload-outlined style="font-size: 8em; color: #52c41a" />
      </div>
      <!-- download -->
      <div class="card">
        <a-cloud-download-outlined style="font-size: 8em; color: #00bcd4" />
      </div>
    </div>

    <div style="margin-top: 1em">
      <p>
        Look over synchronize
        <router-link to="/synchronize-history">History</router-link>
      </p>
    </div>
  </div>
</template>

<script>
import { PageHeader } from "ant-design-vue";
import {
  CloudDownloadOutlined,
  CloudUploadOutlined,
} from "@ant-design/icons-vue";
import { RouterLink } from "vue-router";
export default {
  name: "Synchronize",
  components: {
    RouterLink,
    APageHeader: PageHeader,
    // ACard: Card,
    ACloudDownloadOutlined: CloudDownloadOutlined,
    ACloudUploadOutlined: CloudUploadOutlined,
  },
};
</script>

<style scoped>
.card-container {
  height: 100%;
  width: 100%;
  display: flex;
  justify-content: space-between;
}
.card {
  background: #ffffff;
  max-width: 600px;
  min-width: 300px;
  width: 48%;
  height: 300px;
  display: flex;
  justify-content: center;
  align-items: center;
}
.card:hover {
  cursor: pointer;
  box-shadow: #1890ff 0px 0px 10px;
}
</style>