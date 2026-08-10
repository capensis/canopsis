<template>
  <c-layout class="pl-4" column>
    <c-enabled-field
      :model-value="settings.enabled"
      :label="$t('techMetric.enabled')"
      :disabled="settingsPending"
      @update:model-value="changeTechMetricsEnabled"
    />
    <v-data-table
      :headers="headers"
      :items="metrics"
      :loading="pending"
      loader-height="2"
    >
      <template #no-data>
        <tr>
          <td class="text-center" colspan="2">
            <div>{{ $t('techMetric.noDumps') }}</div>
          </td>
          <td>
            <c-action-btn
              :disabled="exporting"
              :tooltip="$t('techMetric.generateDump')"
              icon="play_circle"
              color="secondary"
              @click="exportTechMetrics"
            />
          </td>
        </tr>
      </template>
      <template #item="{ item }">
        <tr>
          <td
            v-if="!exporting && (getMetricItem(item).disabled || getMetricItem(item).notCreated)"
            class="text-center"
            colspan="2"
          >
            <div>
              {{
                getMetricItem(item).disabled
                  ? $t('techMetric.metricsDisabled')
                  : $t('techMetric.noDumps')
              }}
            </div>
          </td>
          <template v-else>
            <td>{{ formatDate(getMetricItem(item).created) }}</td>
            <td>
              <v-progress-circular v-if="exporting" color="primary" indeterminate />
              <span v-else>
                {{ formatDuration(getMetricItem(item).duration) }}
              </span>
            </td>
          </template>
          <td>
            <c-action-btn
              :disabled="exporting || getMetricItem(item).disabled"
              :tooltip="$t('techMetric.generateDump')"
              icon="play_circle"
              color="secondary"
              @click="exportTechMetrics"
            />
            <c-action-btn
              v-if="getMetricItem(item).readyToDownload"
              :disabled="exporting"
              :tooltip="$t('techMetric.downloadDump')"
              icon="save_alt"
              color="secondary"
              @click="downloadTechMetrics"
            />
          </td>
        </tr>
      </template>
    </v-data-table>
  </c-layout>
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import { TECH_METRICS_EXPORT_STATUSES } from '@/constants';
import { convertDateToString } from '@/helpers/date/date';
import { convertDurationToString } from '@/helpers/date/duration';
import { getTechMetricsDownloadFileUrl } from '@/helpers/entities/metric/url';
import { openUrlInNewTab } from '@/helpers/url';
import { useI18n } from '@/hooks/i18n';
import { useFilePolling } from '@/hooks/polling';
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

    const metrics = computed(() =>
      (metric.value ? [metric.value] : []).map(item => ({
        ...item,
        disabled: item.status === TECH_METRICS_EXPORT_STATUSES.disabled,
        readyToDownload: item.status === TECH_METRICS_EXPORT_STATUSES.success,
        notCreated: item.status === TECH_METRICS_EXPORT_STATUSES.none,
      })),
    );

    const headers = computed(() => [
      { title: t('common.created'), key: 'created', width: '33%', sortable: false },
      { title: t('common.timeTaken'), key: 'duration', width: '33%', sortable: false },
      { title: t('common.actionsLabel'), key: 'actions', width: '33%', sortable: false },
    ]);

    const { pending, handler: fetchTechMetrics } = usePendingHandler(async () => {
      metric.value = await fetchTechMetricsExport();

      return metric.value;
    });

    const { pending: settingsPending, handler: fetchSettings } = usePendingHandler(async () => {
      settings.value = await fetchTechMetricsSettings();
    });

    const { poll: generateFile } = useFilePolling({
      createHandler: createTechMetricsExport,
      fetchHandler: fetchTechMetrics,
      completedStatus: TECH_METRICS_EXPORT_STATUSES.success,
      failedStatus: TECH_METRICS_EXPORT_STATUSES.failed,
    });

    const downloadTechMetrics = () => openUrlInNewTab(getTechMetricsDownloadFileUrl());
    const getMetricItem = item => item.raw ?? item;

    const { pending: exporting, handler: exportTechMetrics } = usePendingHandler(generateFile);

    const changeTechMetricsEnabled = async (enabled) => {
      try {
        settingsPending.value = true;

        settings.value = await updateTechMetricsSettings({
          data: {
            ...settings.value,
            enabled,
          },
        });
      } finally {
        settingsPending.value = false;
      }

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
      getMetricItem,
      formatDate: convertDateToString,
      formatDuration: convertDurationToString,
    };
  },
};
</script>
