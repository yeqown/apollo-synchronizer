<template>
  <div>
    <h2>Apollo Synchronizer</h2>
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
  name: "App",
  components: {
    AStep: Step,
    ASteps: Steps,
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
