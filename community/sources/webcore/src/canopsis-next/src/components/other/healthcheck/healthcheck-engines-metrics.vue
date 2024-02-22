<template>
  <v-data-table
    :headers="headers"
    :items="metrics"
    :loading="pending"
    loader-height="2"
  >
    <template #no-data="">
      <td
        class="text-center"
        colspan="2"
      >
        <div>{{ $t('techMetric.noDumps') }}</div>
      </td>
      <td>
        <c-action-btn
          :disabled="exporting"
          :tooltip="$t('techMetric.generateDump')"
          icon="play_circle_filled"
          color="secondary"
          @click="exportTechMetrics"
        />
      </td>
    </template>
    <template #item="{ item }">
      <td
        v-if="!exporting && (item.disabled || isMetricNotCreated(item))"
        class="text-center"
        colspan="2"
      >
        <div>{{ item.disabled ? $t('techMetric.metricsDisabled') : $t('techMetric.noDumps') }}</div>
      </td>
      <template v-else>
        <td>{{ item.created | date }}</td>
        <td>
          <v-progress-circular
            v-if="exporting"
            color="primary"
            indeterminate
          />
          <span v-else>{{ item.duration | duration }}</span>
        </td>
      </template>
      <td>
        <c-action-btn
          :disabled="exporting || item.disabled"
          :tooltip="$t('techMetric.generateDump')"
          icon="play_circle_filled"
          color="secondary"
          @click="exportTechMetrics"
        />
        <c-action-btn
          v-if="isMetricReadyToDownload(item)"
          :disabled="exporting"
          :tooltip="$t('techMetric.downloadDump')"
          icon="save_alt"
          color="secondary"
          @click="downloadTechMetrics"
        />
      </td>
    </template>
  </v-data-table>
</template>

<script>
import { TECH_METRICS_EXPORT_STATUSES } from '@/constants';

import { getTechMetricsDownloadFileUrl } from '@/helpers/entities/metric/url';

import { entitiesTechMetricsMixin } from '@/mixins/entities/tech-metrics';
import { exportMixinCreator } from '@/mixins/widget/export';

export default {
  mixins: [
    entitiesTechMetricsMixin,
    exportMixinCreator({
      createExport: 'createTechMetricsExport',
      fetchExport: 'fetchTechMetrics',
      completedStatus: TECH_METRICS_EXPORT_STATUSES.success,
      failedStatus: TECH_METRICS_EXPORT_STATUSES.failed,
    }),
  ],
  data() {
    return {
      pending: false,
      exporting: false,
      metric: null,
    };
  },
  computed: {
    headers() {
      return [
        { text: this.$t('common.created'), width: '33%', sortable: false },
        { text: this.$t('common.timeTaken'), width: '33%', sortable: false },
        { text: this.$t('common.actionsLabel'), width: '33%', sortable: false },
      ];
    },

    metrics() {
      const metrics = this.metric ? [this.metric] : [];

      return metrics.map(metric => ({
        ...metric,
        disabled: this.isMetricDisabled(metric),
      }));
    },
  },
  async mounted() {
    this.pending = true;

    await this.fetchTechMetrics();

    this.pending = false;
  },
  methods: {
    async fetchTechMetrics() {
      this.metric = await this.fetchTechMetricsExport();

      return this.metric;
    },

    async exportTechMetrics() {
      this.exporting = true;

      await this.generateFile();

      this.exporting = false;
    },

    async downloadTechMetrics() {
      this.downloadFile(getTechMetricsDownloadFileUrl());
    },

    isMetricReadyToDownload(item) {
      return item.status === TECH_METRICS_EXPORT_STATUSES.success;
    },

    isMetricDisabled(item) {
      return item.status === TECH_METRICS_EXPORT_STATUSES.disabled;
    },

    isMetricNotCreated(item) {
      return item.status === TECH_METRICS_EXPORT_STATUSES.none;
    },
  },
};
</script>
