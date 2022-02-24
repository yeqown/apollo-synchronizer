<template>
  <div>
    <a-page-header
      style="
        border: 1px solid rgb(235, 237, 240);
        margin-bottom: 1em;
        background-color: #ffffff;
      "
      title="Synchronize"
      :sub-title="`${(action || 'unknown').toLowerCase()} namespaces or files.`"
      @back="() => this.$router.push({ name: 'synchronize' })"
    >
      <template #extra>
        <a-cloud-upload-outlined
          v-if="(action || 'unknown').toLowerCase() === 'upload'"
          style="font-size: 2em; color: #52c41a"
        />
        <a-cloud-download-outlined
          v-if="(action || 'unknown').toLowerCase() === 'download'"
          style="font-size: 2em; color: #00bcd4"
        />
      </template>
    </a-page-header>

    <div id="synchronize-form-container">
      <a-row :gutter="24" style="width: 100%">
        <a-col :span="18">
          <h3>Use Setting</h3>
          <a-form
            :model="form"
            :label-col="{ style: { width: '100px' } }"
            :wapper-col="{ span: 14 }"
            style="min-width: 400px"
          >
            <a-form-item label="Apollo">
              <a-select
                v-model:value="form.usingSettingIdx"
                @change="handleSettingChange"
              >
                <template
                  v-for="(item, index) in settings"
                  :key="`setting-${index}`"
                >
                  <a-select-option :value="index">{{
                    `${item.title}(${item.portalAddr})`
                  }}</a-select-option>
                </template>
                <a-select-option
                  :value="-1"
                  :disabled="form.usingSettingIdx !== -1"
                >
                  Choose setting</a-select-option
                >
              </a-select>
            </a-form-item>

            <a-form-item label="Env">
              <a-select
                :disabled="form.usingSettingIdx < 0"
                v-model:value="form.usingEnv"
              >
                <template v-if="form.usingSettingIdx >= 0">
                  <template
                    v-for="(env, index) in settings[form.usingSettingIdx].envs"
                    :key="`env-${index}`"
                  >
                    <a-select-option :value="env">{{ env }}</a-select-option>
                  </template>
                </template>
              </a-select>
            </a-form-item>

            <a-form-item label="Cluster">
              <a-select
                :disabled="form.usingSettingIdx < 0"
                v-model:value="form.usingCluster"
              >
                <template v-if="form.usingSettingIdx >= 0">
                  <template
                    v-for="(cluster, index) in settings[form.usingSettingIdx]
                      .clusters"
                    :key="`cluster-${index}`"
                  >
                    <a-select-option :value="cluster">{{
                      cluster
                    }}</a-select-option>
                  </template>
                </template>
              </a-select>
            </a-form-item>

            <a-form-item label="AppId" name="appId">
              <a-input
                v-model:value="form.appId"
                placeholder="input your apollo appId"
                :rules="[
                  { required: true, message: 'Please input your appId!' },
                ]"
              ></a-input>
            </a-form-item>
            <a-form-item label="Optional">
              <a-checkbox-group
                style="width: 100%"
                v-model:value="form.optional"
              >
                <a-row>
                  <a-col :span="7"
                    ><a-checkbox value="force">Force</a-checkbox></a-col
                  >
                  <a-col :span="9" value="autopublish"
                    ><a-checkbox>AutoPublish</a-checkbox></a-col
                  >
                  <a-col :span="8" value="overwrite"
                    ><a-checkbox>Overwrite</a-checkbox></a-col
                  >
                </a-row>
              </a-checkbox-group>
            </a-form-item>

            <a-form-item :wrapper-col="{ span: 12, offset: 6 }">
              <a-button type="default" @click="() => this.$router.back()"
                >Cancel</a-button
              >
              <a-button
                style="margin-left: 10px"
                type="primary"
                :disabled="form.usingSettingIdx < 0"
                @click="doSynchronize"
                >Synchronize</a-button
              >
            </a-form-item>
          </a-form>
        </a-col>

        <a-col :span="6">
          <h3>Config Overview</h3>
          <div style="background: #f3f4f5; width: 100%; height: 250px"></div>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script>
import {
  PageHeader,
  Form,
  FormItem,
  Select,
  SelectOption,
  Button,
  Row,
  Col,
  CheckboxGroup,
  Checkbox,
  Input,
} from "ant-design-vue";
import {
  CloudDownloadOutlined,
  CloudUploadOutlined,
} from "@ant-design/icons-vue";
import { loadSetting } from "../interact/setting";
import {
  notificationError,
  notificationWarning,
  notificationSuccess,
} from "../utils/notification";
import { modeMapping, containsKey, synchronize } from "../interact/synchronize";
export default {
  name: "SynchronizeDo",
  components: {
    APageHeader: PageHeader,
    AForm: Form,
    AFormItem: FormItem,
    ASelect: Select,
    ASelectOption: SelectOption,
    AButton: Button,
    ARow: Row,
    ACol: Col,
    ACheckboxGroup: CheckboxGroup,
    ACheckbox: Checkbox,
    AInput: Input,
    // icons
    ACloudDownloadOutlined: CloudDownloadOutlined,
    ACloudUploadOutlined: CloudUploadOutlined,
  },
  data: () => ({
    settings: [],
    action: "unknown",
    form: {
      usingSettingIdx: -1,
      usingEnv: "",
      usingCluster: "",
      appId: "",
      optional: [],
    },
  }),
  mounted() {
    this.action = this.$route.params.action;
    loadSetting().then(
      (settings) => {
        this.settings = settings;
      },
      (err) => {
        notificationError(err);
      }
    );
    // TODO(@yeqown): register event listener to recv event and data from background.
  },
  methods: {
    doSynchronize() {
      if (this.form.usingSettingIdx < 0) {
        notificationWarning("Please choose setting first!");
        return;
      }

      if (this.form.usingEnv === "" || this.form.usingCluster === "") {
        notificationWarning("Please choose env and cluster first!");
        return;
      }

      if (this.form.appId === "") {
        notificationWarning("Please input your appId!");
        return;
      }

      const setting = this.settings[this.form.usingSettingIdx];
      let scope = {
        portalAddr: setting.portalAddr,
        account: setting.account,
        fs: setting.fs,
        secret: setting.secret,
        appId: this.form.appId,
        env: this.form.usingEnv,
        cluster: this.form.usingCluster,
        mode: modeMapping[this.action.toLowerCase()],
        isForce: containsKey(this.form.optional, "force"), // FIXME(@yeqown): failed to get value from checkbox
        isOverwrite: containsKey(this.form.optional, "overwrite"), // FIXME(@yeqown): failed to get value from checkbox
        isAutoPublish: containsKey(this.form.optional, "autopublish"), // FIXME(@yeqown): failed to get value from checkbox
      };

      console.log("doSynchronize with scope=", scope);
      synchronize(scope).then(
        () => {
          notificationSuccess("Synchronize finished");
        },
        (err) => {
          notificationError(err);
        }
      );
    },
    handleSettingChange(value) {
      // console.log("handleSettingChange", value);
      this.form.usingSettingIdx = value;
      // reset related fields
      this.form.usingEnv = "";
      this.form.usingCluster = "";
    },
  },
};
</script>

<style scoped>
#synchronize-form-container {
  width: 100%;
  display: flex;
  justify-content: center;

  background: #ffffff;
  height: 400px;
  padding: 16px 24px;
}

/* #synchronize-form-container > h3 {
  text-align: left;
  font-weight: bold;
} */
</style>