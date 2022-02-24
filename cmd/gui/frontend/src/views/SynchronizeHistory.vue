<template>
  <div>
    <a-page-header
      style="
        border: 1px solid rgb(235, 237, 240);
        margin-bottom: 1em;
        background-color: #ffffff;
      "
      title="Synchronize"
      sub-title="display history of sychronization."
      @back="() => this.$router.push({ name: 'synchronize' })"
    />
  </div>
</template>

<script>
import { PageHeader } from "ant-design-vue";
export default {
  name: "SynchronizeHistory",
  components: {
    APageHeader: PageHeader,
  },
};
</script>