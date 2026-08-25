<template>
  <div class="position-relative">
    <c-progress-overlay :pending="remediationStatisticPending" />

    <remediation-statistics-filters v-model="query" :min-date="minDate" />

    <remediation-statistics-chart
      :metrics="metrics"
      :data-type="query.type"
      :sampling="query.sampling"
      :min-date="minDate"
      :downloading="downloading"
      @export:csv="exportRemediationStatisticsAsCsv"
      @export:png="exportRemediationStatisticsAsPng"
    />
  </div>
</template>

<script setup>
import { computed, onMounted, ref, watch } from 'vue';

import {
  EXPORT_FETCHING_INTERVAL,
  REMEDIATION_STATISTICS_FILENAME_PREFIX,
} from '@/config';
import {
  DATETIME_FORMATS,
  EXPORT_STATUSES,
  QUICK_RANGES,
  REMEDIATION_INSTRUCTION_TYPES,
  REMEDIATION_STATISTICS_CHART_DATA_TYPE,
  SAMPLINGS,
} from '@/constants';

import {
  convertDateToStartOfDayTimestampByTimezone,
  convertDateToString,
} from '@/helpers/date/date';
import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';
import { convertMetricsToTimezone } from '@/helpers/entities/metric/list';
import { isMetricsQueryChanged } from '@/helpers/entities/metric/query';
import { getExportMetricDownloadFileUrl } from '@/helpers/entities/metric/url';
import { saveFile } from '@/helpers/file/files';
import { openUrlInNewTab } from '@/helpers/url';

import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useInfo } from '@/hooks/store/modules/info';
import { useMetrics } from '@/hooks/store/modules/metrics';
import { useRemediationStatistic } from '@/hooks/store/modules/remediation-statistic';

import RemediationStatisticsChart from './partials/remediation-statistics-chart.vue';
import RemediationStatisticsFilters from './partials/remediation-statistics-filters.vue';

const { t } = useI18n();
const popups = usePopups();
const { timezone } = useInfo();
const {
  remediationStatisticMetrics,
  remediationStatisticMeta,
  remediationStatisticPending,
  fetchRemediationStatisticList,
} = useRemediationStatistic();
const { createRemediationExport, fetchMetricExport } = useMetrics();

const downloading = ref(false);
const query = ref({
  sampling: SAMPLINGS.day,
  type: REMEDIATION_STATISTICS_CHART_DATA_TYPE.percent,
  instruction: '',
  interval: {
    from: QUICK_RANGES.last7Days.start,
    to: QUICK_RANGES.last7Days.stop,
  },
});

const minDate = computed(() => {
  const value = remediationStatisticMeta.value?.minDate ?? remediationStatisticMeta.value?.min_date;

  return value ? convertDateToStartOfDayTimestampByTimezone(value, timezone.value) : null;
});

const interval = computed(() =>
  convertMetricIntervalToTimestamp({
    interval: query.value.interval,
    timezone: timezone.value,
  }),
);

const metrics = computed(() =>
  convertMetricsToTimezone(
    remediationStatisticMetrics.value.filter(
      ({ timestamp }) => !minDate.value || timestamp >= minDate.value,
    ),
    timezone.value,
  ),
);

const isInstructionType = instruction =>
  [REMEDIATION_INSTRUCTION_TYPES.manual, REMEDIATION_INSTRUCTION_TYPES.auto].includes(instruction);

const getQuery = () => {
  const { instruction } = query.value;
  const result = {
    ...interval.value,

    sampling: query.value.sampling,
  };

  if (isInstructionType(instruction)) {
    result.instruction_type = instruction;
  } else if (instruction) {
    result.instruction = instruction;
  }

  return result;
};

const fetchList = () => fetchRemediationStatisticList({ params: getQuery() });

const getInstructionString = (instruction) => {
  if (!instruction) {
    return t('remediation.statistic.allInstructions');
  }

  return isInstructionType(instruction)
    ? t(`remediation.instruction.types.${instruction}`)
    : instruction;
};

const getFileName = () => {
  const fromTime = convertDateToString(interval.value.from, DATETIME_FORMATS.short);
  const toTime = convertDateToString(interval.value.to, DATETIME_FORMATS.short);
  const instruction = getInstructionString(query.value.instruction);

  return [
    REMEDIATION_STATISTICS_FILENAME_PREFIX,
    fromTime,
    toTime,
    query.value.sampling,
    query.value.type,
    instruction,
  ].join('-');
};

const waitGeneratingFile = id =>
  new Promise((resolve, reject) => {
    setTimeout(async () => {
      try {
        const exportData = await fetchMetricExport({ id });

        if (exportData.status === EXPORT_STATUSES.completed) {
          resolve(exportData);

          return;
        }

        if (exportData.status === EXPORT_STATUSES.failed) {
          reject();

          return;
        }

        resolve(waitGeneratingFile(id));
      } catch (err) {
        reject(err);
      }
    }, EXPORT_FETCHING_INTERVAL);
  });

const exportRemediationStatisticsAsPng = async (blob) => {
  try {
    await saveFile(blob, getFileName());
  } catch (err) {
    console.error(err);

    popups.error({ text: err.message || t('errors.default') });
  }
};

const exportRemediationStatisticsAsCsv = async () => {
  downloading.value = true;

  try {
    const exportRequest = await createRemediationExport({ data: getQuery() });
    const id = exportRequest._id ?? exportRequest.id;
    const fileData = await waitGeneratingFile(id);

    openUrlInNewTab(getExportMetricDownloadFileUrl(fileData._id ?? fileData.id));
  } catch (err) {
    console.error(err);

    popups.error({ text: t('kpi.popups.exportFailed') });
  } finally {
    downloading.value = false;
  }
};

watch(minDate, (value) => {
  if (value && interval.value.from < value) {
    query.value = {
      ...query.value,
      interval: {
        ...query.value.interval,
        from: value,
      },
    };
  }
});

watch(query, (value, oldValue) => {
  if (isMetricsQueryChanged(value, oldValue, minDate.value)) {
    fetchList();
  }
});

onMounted(fetchList);
</script>
