<template lang="pug">
  v-data-table(:headers="headers", :items="metrics", :loading="pending")
    template(#no-data="")
      td.text-xs-center(colspan="2")
        div {{ $t('techMetric.noDumps') }}
      td
        c-action-btn(
          :disabled="exporting",
          :tooltip="$t('techMetric.generateDump')",
          icon="play_circle_filled",
          color="secondary",
          @click="exportTechMetrics"
        )
    template(#items="{ item }")
      td.text-xs-center(
        v-if="!exporting && (item.disabled || isMetricNotCreated(item))",
        colspan="2"
      )
        div {{ item.disabled ? $t('techMetric.metricsDisabled') : $t('techMetric.noDumps') }}
      template(v-else)
        td {{ item.created | date }}
        td
          v-progress-circular(v-if="exporting", color="primary", indeterminate)
          span(v-else) {{ item.duration | duration }}
      td
        c-action-btn(
          :disabled="exporting || item.disabled",
          :tooltip="$t('techMetric.generateDump')",
          icon="play_circle_filled",
          color="secondary",
          @click="exportTechMetrics"
        )
        c-action-btn(
          v-if="isMetricReadyToDownload(item)",
          :disabled="exporting",
          :tooltip="$t('techMetric.downloadDump')",
          icon="save_alt",
          color="secondary",
          @click="downloadTechMetrics"
        )
</template>

<script>
import { API_HOST, API_ROUTES } from '@/config';

import { TECH_METRICS_EXPORT_STATUSES } from '@/constants';

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
      this.downloadFile(`${API_HOST}${API_ROUTES.techMetrics}/download`);
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
