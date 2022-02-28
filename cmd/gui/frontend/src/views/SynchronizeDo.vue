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
                    ><a-checkbox value="force"
                      >Force
                      <a-popover>
                        <template #content>
                          <p>
                            If not set, creation and deletion will be ignored!
                          </p>
                        </template>
                        <a-question-circle-two-tone />
                      </a-popover>
                    </a-checkbox>
                  </a-col>
                  <a-col :span="8">
                    <a-checkbox value="overwrite"
                      >Overwrite
                      <a-popover>
                        <template #content>
                          <p>If not set, modification will be ignored!</p>
                        </template>
                        <a-question-circle-two-tone />
                      </a-popover>
                    </a-checkbox>
                  </a-col>
                  <a-col :span="9">
                    <a-checkbox
                      value="autopublish"
                      :disabled="action !== 'upload'"
                      >AutoPublish
                      <a-popover>
                        <template #content>
                          <p>
                            If not set, changes those uploaded to apollo would
                            not be published automatically!
                          </p>
                        </template>
                        <a-question-circle-two-tone />
                      </a-popover>
                    </a-checkbox>
                  </a-col>
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

    <a-modal
      v-model:visible="syncModalVisible"
      width="100%"
      :closable="true"
      :footer="null"
      :afterClose="handleModalCloseCallback"
    >
      <div id="synchronize-render" style="width: 100%">
        <a-steps :current="syncStepsCurrent">
          <a-step title="Fetch">
            ><template #icon><a-file-sync-outlined /> </template>
          </a-step>
          <a-step title="Decide"
            ><template #icon> <a-code-outlined /></template
          ></a-step>
          <a-step title="Execute"
            ><template #icon> <a-cloud-sync-outlined /> </template
          ></a-step>
          <a-step title="Result"
            ><template #icon> <a-notification-outlined /> </template
          ></a-step>
        </a-steps>

        <div id="step-content">
          <div>
            <a-spin tip="executing, wait please" :spinning="loading">
              <template v-if="syncStepsCurrent === 0">
                <a-result title="Fetching namespaces from apollo, please wait!">
                  <template #icon>
                    <a-sync-outlined spin />
                  </template>
                </a-result>
              </template>

              <template v-if="syncStepsCurrent === 1">
                <a-table
                  :columns="[
                    { title: 'Namespace', dataIndex: 'key', key: 'key' },
                    {
                      title: 'Operation',
                      dataIndex: 'mode',
                      key: 'mode',
                      slots: { customRender: 'operation' },
                    },
                    {
                      title: 'File',
                      dataIndex: 'absFilepath',
                      key: 'absFilepath',
                    },
                  ]"
                  :pagination="{ hideOnSinglePage: true }"
                  :data-source="syncRenderDiffs"
                  size="small"
                  bordered
                  :scroll="{ y: 250 }"
                >
                  <template #operation="{ text }">
                    <!-- render yellow text when text is 'M~', green text when text is 'C+', red text when text is 'D-' -->
                    <span v-if="text === 'M~'">
                      <a-tag color="orange">M~</a-tag>
                    </span>
                    <span v-if="text === 'C+'">
                      <a-tag color="green">C+</a-tag>
                    </span>
                    <span v-if="text === 'D-'">
                      <a-tag color="#f50">D-</a-tag>
                    </span>
                  </template>
                </a-table>

                <div class="step-content-footer">
                  <a-button
                    type="primary"
                    @click="confirmSynchronize"
                    :disabled="!syncContinueAble"
                    :loading="loading"
                    >Confirm</a-button
                  >
                </div>
              </template>

              <template v-if="syncStepsCurrent === 2">
                <a-table
                  :columns="[
                    { title: 'Namespace', dataIndex: 'key', key: 'key' },
                    {
                      title: 'Operation',
                      dataIndex: 'mode',
                      key: 'mode',
                      slots: { customRender: 'operation' },
                    },
                    {
                      title: 'Executed Result',
                      dataIndex: 'error',
                      key: 'result',
                      slots: { customRender: 'result' },
                    },
                    {
                      title: 'Published',
                      dataIndex: 'published',
                      key: 'published',
                      slots: { customRender: 'published' },
                    },
                  ]"
                  :pagination="{ hideOnSinglePage: true }"
                  :data-source="syncRenderResults"
                  size="small"
                  bordered
                  :scroll="{ y: 250 }"
                >
                  <template #operation="{ text }">
                    <!-- render yellow text when text is 'M~', green text when text is 'C+', red text when text is 'D-' -->
                    <span v-if="text === 'M~'">
                      <a-tag color="orange">M~</a-tag>
                    </span>
                    <span v-if="text === 'C+'">
                      <a-tag color="green">C+</a-tag>
                    </span>
                    <span v-if="text === 'D-'">
                      <a-tag color="#f50">D-</a-tag>
                    </span>
                  </template>

                  <template #result="{ text, record }">
                    <a-badge
                      dot
                      :color="record.succeeded ? 'green' : 'orange'"
                      :text="record.succeeded ? '' : text"
                    />
                  </template>

                  <template #published="{ record }">
                    <!-- render red badge if text.published is true otherwire render green badge -->
                    <a-badge
                      dot
                      :color="record.published ? 'green' : 'orange'"
                    />
                  </template>
                </a-table>

                <div class="step-content-footer">
                  <a-button
                    type="primary"
                    @click="
                      () => {
                        syncStepsCurrent = 3;
                      }
                    "
                    :disabled="!syncContinueAble"
                    :loading="loading"
                    >Continue</a-button
                  >
                </div>
              </template>

              <!-- setp 4 -->
              <template v-if="syncStepsCurrent === 3">
                <a-result
                  :status="syncResult.success ? 'success' : 'error'"
                  :title="
                    syncResult.success
                      ? 'Synchronize Successfully!'
                      : syncResult.failedReason
                  "
                >
                </a-result>
              </template>
            </a-spin>
          </div>

          <!-- button for debug -->
          <!-- <a-button
            type="primary"
            @click="() => (syncStepsCurrent = ++syncStepsCurrent % 4)"
            >Debug Next</a-button
          > -->
        </div>
      </div>
    </a-modal>
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
  Popover,
  Modal,
  Steps,
  Step,
  Result,
  Table,
  Tag,
  Badge,
  Spin,
} from "ant-design-vue";
import {
  CloudDownloadOutlined,
  CloudUploadOutlined,
  QuestionCircleTwoTone,
  SyncOutlined,
  NotificationOutlined,
  CodeOutlined,
  CloudSyncOutlined,
  FileSyncOutlined,
} from "@ant-design/icons-vue";
import { loadSetting } from "../interact/setting";
import {
  notificationError,
  notificationWarning,
  notificationSuccess,
} from "../utils/notification";
import {
  modeMapping,
  containsKey,
  synchronize,
  bindEventOnce,
  unbindEvent,
  EVENT_RENDER_DIFF,
  EVENT_RENDER_RESULT,
  inputDecide,
  decideMapping,
} from "../interact/synchronize";
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
    APopover: Popover,
    AModal: Modal,
    ASteps: Steps,
    AStep: Step,
    AResult: Result,
    ATable: Table,
    ATag: Tag,
    ABadge: Badge,
    ASpin: Spin,
    // icons
    ACloudDownloadOutlined: CloudDownloadOutlined,
    ACloudUploadOutlined: CloudUploadOutlined,
    AQuestionCircleTwoTone: QuestionCircleTwoTone,
    ASyncOutlined: SyncOutlined,
    ACloudSyncOutlined: CloudSyncOutlined,
    ANotificationOutlined: NotificationOutlined,
    ACodeOutlined: CodeOutlined,
    AFileSyncOutlined: FileSyncOutlined,
  },
  data: () => ({
    loading: false,
    settings: [],
    action: "unknown",
    form: {
      usingSettingIdx: -1,
      usingEnv: "",
      usingCluster: "",
      appId: "",
      optional: ["force", "overwrite"],
    },
    syncRenderDiffs: [
      // {
      //   key: "mock.yaml",
      //   mode: "C+",
      //   absFilepath: "path/to",
      // },
    ],
    syncContinueAble: true,
    syncRenderResults: [
      // {
      //   key: "mock.yaml",
      //   mode: "M~",
      //   error: "hahah",
      //   published: true,
      //   succeeded: true,
      // },
    ],
    syncStepsCurrent: 0,
    syncModalVisible: false,
    syncResult: {
      success: false,
      failedReason: "",
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

    // DONE(@yeqown): register event listener to recv event and data from background.
    bindEventOnce(EVENT_RENDER_DIFF, (data) => {
      console.log("event render diff triggered", data);
      this.syncRenderDiffs = data;
      setTimeout(() => {
        this.loading = false;
        this.syncStepsCurrent = 1;
      }, 1000);
    });

    bindEventOnce(EVENT_RENDER_RESULT, (data) => {
      console.log("event render result triggered", data);
      this.syncRenderResults = data;
      setTimeout(() => {
        this.loading = false;
        this.syncStepsCurrent = 2;
      }, 500);
    });
  },
  unmounted() {
    unbindEvent(EVENT_RENDER_DIFF);
    unbindEvent(EVENT_RENDER_RESULT);
  },
  methods: {
    doSynchronize() {
      console.log("============ form values", this.form);
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

      if (
        !this.action ||
        (this.action !== "upload" && this.action !== "download")
      ) {
        notificationWarning("Unknown action!, please re-enter this page");
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
        isForce: containsKey(this.form.optional, "force"),
        isOverwrite: containsKey(this.form.optional, "overwrite"),
        isAutoPublish: containsKey(this.form.optional, "autopublish"),
      };

      const quickFail = (reason) => {
        this.syncStepsCurrent = 3;
        this.syncResult = {
          success: false,
          failedReason: reason,
        };
      };

      console.log("doSynchronize with scope=", scope);
      synchronize(scope).then(
        (result) => {
          if (result.succeeded) {
            notificationSuccess("Synchronize succeed!");
            this.syncResult = {
              success: true,
              failedReason: "",
            };
          } else {
            notificationError("Synchronize failed: " + result.failedReason);
            quickFail(result.failedReason);
          }
        },
        (err) => {
          notificationError(err);
          this.syncContinueAble = false;
          quickFail(err);
        }
      );

      this.syncModalVisible = true;
    },
    confirmSynchronize() {
      console.log("confirmSynchronize============");
      inputDecide(decideMapping["confirm"]);
      this.loading = true;
    },
    handleSettingChange(value) {
      // console.log("handleSettingChange", value);
      this.form.usingSettingIdx = value;
      // reset related fields
      this.form.usingEnv = "";
      this.form.usingCluster = "";
    },
    handleModalCloseCallback() {
      console.log("handleModalCloseCallback called");
      // modal close callback
      this.loading = false;
      this.syncStepsCurrent = 0;
      this.loading = false;
      this.syncRenderDiffs = [];
      this.syncRenderResults = [];
      this.syncResult = { success: false, failedReason: "unset" };
    },
  },
};
</script>

<style scoped>
#synchronize-form-container {
  width: 100%;
  max-width: 800px;
  display: flex;
  justify-content: center;

  background: #ffffff;
  min-height: 400px;
  padding: 16px 24px;
}

#synchronize-render {
  margin-top: 1.5em;
  max-height: 450px;
}

#step-content {
  /* overflow: scroll; */
  height: 300px;
  width: 100%;
  margin-top: 1em;
}

.step-content-footer {
  display: flex;
  margin-top: 1em;
  justify-content: center;
}
</style>