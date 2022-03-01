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
        <a-button type="primary" @click="handleOpenEditModal"> Add </a-button>
        <a-button type="default" :disabled="!modified" @click="save">
          Save
        </a-button>
      </template>
    </a-page-header>

    <!-- setting render -->
    <div
      style="
        height: 400px;
        text-align: left;
        background: #ffffff;
        overflow: scroll;
      "
    >
      <!-- empty -->
      <div
        v-if="!settings || settings.length === 0"
        style="
          height: 100%;
          display: flex;
          justify-content: center;
          align-items: center;
        "
      >
        <a-empty
          image="https://gw.alipayobjects.com/mdn/miniapp_social/afts/img/A*pevERLJC9v0AAAAAAAAAAABjAQAAAQ/original"
          :image-style="{
            height: '60px',
          }"
        >
          <template #description>
            <span> How about adding a apollo cluster?</span>
          </template>
          <a-button type="primary" @click="handleOpenEditModal"
            >Add Now</a-button
          >
        </a-empty>
      </div>

      <!-- list -->
      <a-list
        v-else
        item-layout="vertical"
        size="small"
        :data-source="settings"
      >
        <template #renderItem="{ item, index }">
          <a-list-item :key="item.title">
            <template #actions>
              <DeleteTwoTone
                two-tone-color="#eb2f96"
                @click="handleDeleteSetting(index)"
              />
              <EditTwoTone
                two-tone-color="#1890ff"
                @click="handleEditSetting(index)"
              />
            </template>

            <!-- config data -->
            <a-descriptions
              bordered
              size="small"
              :column="{ xxl: 4, xl: 3, lg: 3, md: 3, sm: 2, xs: 1 }"
            >
              <a-descriptions-item
                label="Title"
                :labelStyle="{
                  backgroundColor: '#f5f5f5',
                  fontWeight: 'bold',
                }"
              >
                {{ item.title }}
              </a-descriptions-item>
              <a-descriptions-item
                label="Portal"
                :span="2"
                :labelStyle="{
                  backgroundColor: '#f5f5f5',
                  fontWeight: 'bold',
                }"
                :contentStyle="{ color: '#1a90ff' }"
              >
                <a @click="_openURL(item.portalAddr)"> {{ item.portalAddr }}</a>
              </a-descriptions-item>
              <a-descriptions-item
                label="Apollo Config"
                :span="3"
                :labelStyle="{
                  backgroundColor: '#f5f5f5',
                  fontWeight: 'bold',
                }"
              >
                Account: {{ item.account }}
                <br />
                Envs: {{ item.envs }}
                <br />
                Clusters: {{ item.clusters }}
                <br />
                Secret: {{ item.secret }}
              </a-descriptions-item>
              <a-descriptions-item
                label="Directory"
                :span="3"
                :labelStyle="{
                  backgroundColor: '#f5f5f5',
                  fontWeight: 'bold',
                }"
              >
                <FolderOpenOutlined
                  @click="_openURL(item.fs)"
                  style="margin-right: 1em"
                />
                <a @click="_openURL(item.fs)">{{ item.fs }}</a>
              </a-descriptions-item>
            </a-descriptions>
          </a-list-item>
        </template>
      </a-list>
    </div>

    <a-modal
      :visible="editModalVisible"
      title="Edit setting"
      wrap-class-name="full-modal"
      width="100%"
      @ok="handleEditModalOk"
      @cancel="
        () => {
          editModalVisible = false;
          modified = false;

          this.resetForm();
        }
      "
    >
      <a-form
        :model="form"
        :label-col="{ style: { width: '140px' } }"
        :wrapper-col="{ span: 12 }"
      >
        <a-form-item
          label="Title"
          name="title"
          :rules="[
            { required: true, message: 'Please input your setting title!' },
          ]"
        >
          <a-input
            v-model:value="form.title"
            placeholder="input title of current setting"
          >
            <template #prefix>
              <FieldStringOutlined class="site-form-item-icon" />
            </template>
          </a-input>
        </a-form-item>
        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="Apollo Account">
              <a-input v-model:value="form.account" placeholder="apollo">
                <template #prefix>
                  <UserOutlined class="site-form-item-icon" />
                </template>
              </a-input>
            </a-form-item>
          </a-col>

          <a-col :span="12">
            <a-form-item label="Apollo Secret">
              <a-input v-model:value="form.secret" placeholder="applied token">
                <template #prefix>
                  <LockOutlined class="site-form-item-icon" />
                </template>
              </a-input>
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="Apollo Portal">
              <a-input
                v-model:value="form.portalAddr"
                placeholder="http://example.com"
              >
                <template #prefix>
                  <IeOutlined class="site-form-item-icon" />
                </template>
              </a-input>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="Local Directory">
              <a-input
                v-model:value="form.fs"
                placeholder="the path to save remote namespace"
              >
                <template #prefix>
                  <FolderOpenOutlined class="site-form-item-icon" />
                </template>
              </a-input>
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="Apollo Envs">
              <a-input
                v-model:value="formAddEnv"
                type="text"
                @blur="handleEnvAddConfirm"
                @keyup.enter="handleEnvAddConfirm"
              >
                <template #prefix>
                  <EnvironmentOutlined class="site-form-item-icon" />
                </template>
              </a-input>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="Values">
              <template v-for="tag in form.envs" :key="tag">
                <a-tag
                  color="green"
                  :closable="true"
                  @close="handleEnvCloseTag(tag)"
                >
                  {{ tag }}
                </a-tag>
              </template>
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="24">
          <a-col :span="12">
            <a-form-item label="Apollo Clusters">
              <a-input
                v-model:value="formAddCluster"
                type="text"
                @blur="handleClusterAddConfirm"
                @keyup.enter="handleClusterAddConfirm"
              >
                <template #prefix>
                  <ClusterOutlined class="site-form-item-icon" />
                </template>
              </a-input>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="Values">
              <template v-for="tag in form.clusters" :key="tag">
                <a-tag
                  color="blue"
                  :closable="true"
                  @close="handleClusterCloseTag(tag)"
                >
                  {{ tag }}
                </a-tag>
              </template>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-modal>
  </div>
</template>

<script>
// import { ref } from "vue";
import { loadSetting, saveSetting } from "../interact/index";
import {
  notificationError,
  notificationSuccess,
  notificationWarning,
} from "../utils/notification";
import {
  PageHeader,
  Button,
  List,
  ListItem,
  Empty,
  // Avatar,
  Descriptions,
  DescriptionsItem,
  Modal,
  Form,
  FormItem,
  Input,
  Tag,
  Row,
  Col,
} from "ant-design-vue";
import {
  UserOutlined,
  LockOutlined,
  IeOutlined,
  FieldStringOutlined,
  FolderOpenOutlined,
  ClusterOutlined,
  EnvironmentOutlined,
  DeleteTwoTone,
  EditTwoTone,
} from "@ant-design/icons-vue";

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
    AModal: Modal,
    AForm: Form,
    AFormItem: FormItem,
    AInput: Input,
    ATag: Tag,
    ARow: Row,
    ACol: Col,
    // icons
    UserOutlined,
    LockOutlined,
    IeOutlined,
    FieldStringOutlined,
    FolderOpenOutlined,
    ClusterOutlined,
    EnvironmentOutlined,
    DeleteTwoTone,
    EditTwoTone,
  },
  data() {
    return {
      modified: false,
      editModalVisible: false,
      settings: [],
      form: {},
      formAddEnv: "",
      formAddCluster: "",
      editingIndex: -1,
    };
  },
  mounted() {
    // console.log("mounted");
    loadSetting().then(
      (settings) => {
        console.log("loadSetting", settings);
        this.settings = settings;
      },
      (error) => {
        notificationError(error);
      }
    );

    this.resetForm();
  },
  methods: {
    resetForm() {
      this.form = {
        title: "",
        account: "apollo",
        clusters: ["default"],
        envs: ["DEV"],
        portalAddr: "https://",
        secret: "",
        fs: "",
      };
    },
    handleOpenEditModal() {
      this.editModalVisible = true;
    },
    handleEditModalOk() {
      if (this.editingIndex >= 0) {
        this.settings[this.editingIndex] = this.form;
      } else {
        this.settings.push(this.form);
      }

      this.editModalVisible = false;
      this.modified = true;
    },
    save() {
      saveSetting(this.settings).then(
        () => {
          // console.log(this.settings);
          console.log("SaveSettings called: ");
          notificationSuccess("All settings have been saved.");
          this.modified = false;
        },
        (error) => {
          console.error(error);
          notificationError(error);
          // DO NOT change modified status
        }
      );
    },
    handleDeleteSetting(index) {
      let settings = this.settings;
      settings.splice(index, 1);
      this.settings = settings;
    },
    handleEditSetting(index) {
      console.log("handleEditSetting index=", index);
      if (index >= 0) {
        this.editingIndex = index;
        this.form = this.settings[index];
        this.editModalVisible = true;
      } else {
        notificationWarning("Please select a setting to edit.");
      }
    },
    handleEnvAddConfirm() {
      // console.log("handleEnvAddConfirm", this.formAddEnv);
      const env = this.formAddEnv;
      let envs = this.form.envs;
      if (env && envs.indexOf(env) === -1) {
        envs = [...envs, env];
      }
      this.form.envs = envs;
      // console.log("handleEnvAddConfirm", this.form.envs);
      this.formAddEnv = "";
    },
    handleEnvCloseTag(removeEnv) {
      let envs = this.form.envs.filter((env) => env !== removeEnv);
      this.form.envs = envs;
    },
    handleClusterAddConfirm() {
      const cluster = this.formAddCluster;
      let clusters = this.form.clusters;
      if (cluster && clusters.indexOf(cluster) === -1) {
        clusters = [...clusters, cluster];
      }
      this.form.clusters = clusters;
      // console.log("clusters", this.form.clusters);
      this.formAddCluster = "";
    },
    handleClusterCloseTag(removeCluster) {
      let clusters = this.form.clusters.filter(
        (cluster) => cluster !== removeCluster
      );
      this.form.clusters = clusters;
    },
    _openURL(url) {
      if (window.runtime && window.runtime.BrowserOpenURL) {
        window.runtime.BrowserOpenURL(url);
        return;
      }

      console.warn(`window.runtime.BrowserOpenURL is not supported.`);
    },
  },
};
</script>

<style scoped>
.full-modal > .ant-modal {
  max-width: 100%;
  top: 0;
  padding-bottom: 0;
  margin: 0;
}

.full-modal > .ant-modal-content {
  display: flex;
  flex-direction: column;
  height: calc(100vh);
}

.full-modal > .ant-modal-body {
  flex: 1;
}

.ant-list-vertical .ant-list-item-extra {
  margin-left: 12px;
}
</style>