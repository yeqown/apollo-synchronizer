<template>
  <div>
    <a-page-header
      title="Dashboard"
      sub-title="overview of the software, and the status of the system."
      style="
        border: 1px solid rgb(235, 237, 240);
        margin-bottom: 1em;
        background-color: #ffffff;
      "
    ></a-page-header>

    <!-- dashboard content -->
    <div style="text-align: left; background: #fff; padding: 16px 24px">
      <!-- usage section -->
      <div>
        <h3>Usage</h3>
        <a-statistic
          title="Active Users"
          :value="112893"
          style="margin-right: 50px"
        />
        <a-statistic
          title="Account Balance (CNY)"
          :precision="2"
          :value="112893"
        />
      </div>
      <!-- document section -->
      <div>
        <h3>Document</h3>
        <p>https://github.com/yeqown/apollo-synchronizer</p>
      </div>
    </div>
    <!-- dashboard content end here -->
  </div>
</template>

<script>
import { Statistic, PageHeader } from "ant-design-vue";
export default {
  name: "Dashboard",
  components: {
    APageHeader: PageHeader,
    AStatistic: Statistic,
  },
};
</script>