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
    <div
      style="
        text-align: left;
        background: #fff;
        padding: 16px 24px;
        height: 400px;
      "
    >
      <!-- statistics section -->
      <div>
        <h3 style="color: #03a9f4; font-weight: bold; font-size: 1.4em">
          #Statistic#
        </h3>
        <div class="statistics-container">
          <a-statistic
            title="Last Open"
            :value="_formatTs(statistics.lastOpenTs)"
            class="statistics-item"
          />
          <a-statistic
            title="First Open"
            :value="_formatTs(statistics.firstOpenTs)"
            class="statistics-item"
          />
          <a-statistic
            title="Open Count"
            :value="_formatNumber(statistics.openCount)"
            class="statistics-item"
          />
          <a-statistic
            title="Open Time"
            :value="_humanizeTime(statistics.openTime)"
            class="statistics-item"
          />
        </div>

        <!-- upload -->
        <div class="statistics-container">
          <a-statistic
            title="Upload Total"
            :value="_formatNumber(statistics.uploadFailedCount)"
            class="statistics-item"
          >
            <template #suffix>
              <span> {{ `/ ${_formatNumber(statistics.uploadCount)}` }} </span>
            </template>
          </a-statistic>
          <a-statistic
            title="Upload Files"
            :value="_formatNumber(statistics.uploadFileCount)"
            class="statistics-item"
          />
          <a-statistic
            title="Upload File Bytes"
            :value="_formatNumber(statistics.uploadFileSize)"
            class="statistics-item"
          />
        </div>

        <!-- download -->
        <div class="statistics-container">
          <a-statistic
            title="Download Total"
            :value="_formatNumber(statistics.downloadFailedCount)"
            class="statistics-item"
          >
            <template #suffix>
              <span>
                {{ `/ ${_formatNumber(statistics.downloadCount)}` }}
              </span>
            </template>
          </a-statistic>
          <a-statistic
            title="Download Files"
            :value="_formatNumber(statistics.downloadFileCount)"
            class="statistics-item"
          />
          <a-statistic
            title="Download File Bytes"
            :value="_formatNumber(statistics.downloadFileSize)"
            class="statistics-item"
          />
        </div>
      </div>

      <!-- document section -->
      <div style="margin-top: 1em">
        <h3 style="color: #03a9f4; font-weight: bold; font-size: 1.4em">
          #Document#
        </h3>
        <a>https://github.com/yeqown/apollo-synchronizer</a>
      </div>
    </div>
    <!-- dashboard content end here -->
  </div>
</template>

<script>
import { loadStatistics } from "../interact/index";
import { formatNumber } from "../utils/index";
import { formatTs, humanizeTime } from "../utils/time";
import { Statistic, PageHeader } from "ant-design-vue";
import { notificationError } from "../utils/notification";

export default {
  name: "Dashboard",
  components: {
    APageHeader: PageHeader,
    AStatistic: Statistic,
  },
  data() {
    return {
      statistics: {
        lastOpenTs: 0,
        firstOpenTs: 0,
        openCount: 0,
        openTime: 0,

        uploadCount: 0,
        uploadFileCount: 0,
        uploadFileSize: 0,
        uploadFailedCount: 0,

        downloadCount: 0,
        downloadFileCount: 0,
        downloadFileSize: 0,
        downloadFailedCount: 0,
      },
    };
  },
  methods: {
    _humanizeTime(ts) {
      return humanizeTime(ts);
    },
    _formatTs(ts) {
      return formatTs(ts);
    },
    _formatNumber(num) {
      return formatNumber(num);
    },
  },
  mounted() {
    loadStatistics().then(
      (statistic) => {
        this.statistics = statistic;
      },
      (error) => {
        notificationError(error);
      }
    );
  },
};
</script>

<style scoped>
.statistics-container {
  display: flex;
  flex-direction: row;
  align-content: flex-start;
}
.statistics-item {
  width: 181px;
  height: 64px;
  overflow: hidden;
}

.upload-text {
  color: #52c41a;
}

.download-text {
  color: #00bcd4;
}
</style>