<template>
  <div class="position-relative">
    <remediation-instruction-stats-list
      :remediation-instruction-stats="remediationInstructionStats"
      :pending="remediationInstructionStatsPending"
      :options.sync="options"
      :total-items="remediationInstructionStatsMeta.total_count"
      :accumulated-before="remediationInstructionStatsMeta.accumulated_before"
      :interval="queryInterval"
      @rate="showRateInstructionModal"
    />
  </div>
</template>

<script>
import { computed, onMounted } from 'vue';

import { QUICK_RANGES } from '@/constants';

import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';

import { useFetchListWithOptions } from '@/hooks/query/shared';
import { useRemdeitionInstructionStats } from '@/hooks/store/modules/remediation-instruction-stats';

import RemediationInstructionStatsList from '@/components/other/remediation/instruction-stats/remediation-instruction-stats-list.vue';

import { useRemediationInstructionStatsRate } from './hooks/remediation-instruction-stats';

export default {
  components: {
    RemediationInstructionStatsList,
  },
  setup() {
    const {
      remediationInstructionStats,
      remediationInstructionStatsPending,
      remediationInstructionStatsMeta,
      fetchRemediationInstructionStatsList,
    } = useRemdeitionInstructionStats();

    const {
      options,
      updateOptions,
      handler: fetchList,
    } = useFetchListWithOptions({
      initialQuery: {
        page: 1,
        itemsPerPage: 10,
        interval: {
          from: QUICK_RANGES.last7Days.start,
          to: QUICK_RANGES.last7Days.stop,
        },
      },
      fetchListHandler: ({ params }) => fetchRemediationInstructionStatsList({
        params: {
          ...params,
          with_flags: true,
        },
      }),
    });

    const queryInterval = computed(() => (
      convertMetricIntervalToTimestamp({ interval: options.value.interval })
    ));

    const { showRateInstructionModal } = useRemediationInstructionStatsRate(fetchList);

    onMounted(fetchList);

    return {
      remediationInstructionStats,
      remediationInstructionStatsPending,
      remediationInstructionStatsMeta,
      options,
      updateOptions,
      queryInterval,
      fetchList,
      showRateInstructionModal,
    };
  },
};
</script>
