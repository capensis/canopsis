<template>
  <v-layout class="pl-4" column>
    <c-enabled-field
      :value="settings.enabled"
      :label="$t('techMetric.enabled')"
      :disabled="settingsPending"
      @change="changeTechMetricsEnabled"
    />
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
          v-if="!exporting && (item.disabled || item.notCreated)"
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
            v-if="item.readyToDownload"
            :disabled="exporting"
            :tooltip="$t('techMetric.downloadDump')"
            icon="save_alt"
            color="secondary"
            @click="downloadTechMetrics"
          />
        </td>
      </template>
    </v-data-table>
  </v-layout>
</template>

<script>
import { ref, computed, onMounted } from 'vue';

import { TECH_METRICS_EXPORT_STATUSES } from '@/constants';

import { getTechMetricsDownloadFileUrl } from '@/helpers/entities/metric/url';

import { useI18n } from '@/hooks/i18n';
import { useExportFile } from '@/hooks/export-file';
import { usePendingHandler } from '@/hooks/query/pending';
import { useTechMetrics } from '@/hooks/store/modules/tech-metrics';

export default {
  setup() {
    const settings = ref({ enabled: false });
    const metric = ref(null);

    const { t } = useI18n();
    const {
      createTechMetricsExport,
      fetchTechMetricsExport,
      fetchTechMetricsSettings,
      updateTechMetricsSettings,
    } = useTechMetrics();

    const metrics = computed(() => (
      (metric.value ? [metric.value] : []).map(item => ({
        ...item,

        disabled: item.status === TECH_METRICS_EXPORT_STATUSES.disabled,
        readyToDownload: item.status === TECH_METRICS_EXPORT_STATUSES.success,
        notCreated: item.status === TECH_METRICS_EXPORT_STATUSES.none,
      }))
    ));

    const headers = computed(() => [
      { text: t('common.created'), width: '33%', sortable: false },
      { text: t('common.timeTaken'), width: '33%', sortable: false },
      { text: t('common.actionsLabel'), width: '33%', sortable: false },
    ]);

    /**
     * Handles fetching and managing the state of technical metrics.
     */
    const {
      pending,
      handler: fetchTechMetrics,
    } = usePendingHandler(async () => {
      metric.value = await fetchTechMetricsExport();

      return metric.value;
    });

    /**
     * Handles fetching and managing the state of settings.
     */
    const {
      pending: settingsPending,
      handler: fetchSettings,
    } = usePendingHandler(async () => {
      settings.value = await fetchTechMetricsSettings();
    });

    /**
     * Handles file generation and download for technical metrics export.
     */
    const { generateFile, downloadFile } = useExportFile({
      createHandler: createTechMetricsExport,
      fetchHandler: fetchTechMetrics,
      completedStatus: TECH_METRICS_EXPORT_STATUSES.success,
      failedStatus: TECH_METRICS_EXPORT_STATUSES.failed,
    });

    /**
     * Initiates the download of the technical metrics file.
     */
    const downloadTechMetrics = () => downloadFile(getTechMetricsDownloadFileUrl());

    /**
     * Handles exporting technical metrics and managing the state of the export process.
     */
    const {
      pending: exporting,
      handler: exportTechMetrics,
    } = usePendingHandler(generateFile);

    /**
     * Updates the enabled status of technical metrics settings.
     *
     * @param {boolean} enabled - The new enabled status.
     */
    const changeTechMetricsEnabled = async (enabled) => {
      settingsPending.value = true;

      settings.value = await updateTechMetricsSettings({
        data: {
          ...settings.value,
          enabled,
        },
      });

      settingsPending.value = false;

      return fetchTechMetrics();
    };

    onMounted(() => {
      fetchTechMetrics();
      fetchSettings();
    });

    return {
      settings,

      metrics,
      headers,

      pending,
      settingsPending,
      exporting,

      exportTechMetrics,
      downloadTechMetrics,
      changeTechMetricsEnabled,
    };
  },
};
</script>
