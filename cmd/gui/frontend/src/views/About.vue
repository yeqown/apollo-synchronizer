<template>
  <div>
    <a-steps :current="1">
      <a-step>
        <!-- <span slot="title">Finished</span> -->
        <template #title>Finished</template>
        <template #description>
          <span>This is a description.</span>
        </template>
      </a-step>
      <a-step title="In Progress" description="This is a description." />
      <a-step title="Waiting" description="This is a description." />
    </a-steps>
  </div>
</template>

<script>
import { Step, Steps } from "ant-design-vue";
export default {
  name: "About",
  components: {
    AStep: Step,
    ASteps: Steps,
  },
};
</script>