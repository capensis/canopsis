<template lang="pug">
  div
    c-page-header
    remediation-instruction-stats-list(
      :remediation-instruction-stats="remediationInstructionStats",
      :pending="remediationInstructionStatsPending",
      :pagination.sync="pagination",
      :total-items="remediationInstructionStatsMeta.total_count",
      :accumulated-before="remediationInstructionStatsMeta.accumulated_before",
      @rate="showRateInstructionModal"
    )
    c-fab-btn(@refresh="fetchList")
</template>

<script>
import {
  MODALS, DATETIME_FORMATS, DATETIME_INTERVAL_TYPES, QUICK_RANGES,
} from '@/constants';

import { dateParse } from '@/helpers/date/date-intervals';
import { convertDateToTimestampByTimezone } from '@/helpers/date/date';

import { authMixin } from '@/mixins/auth';
import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesRemediationInstructionStatsMixin } from '@/mixins/entities/remediation/instruction-stats';
import { entitiesRemediationInstructionsMixin } from '@/mixins/entities/remediation/instructions';

import RemediationInstructionStatsList from '@/components/other/remediation/instruction-stats/remediation-instruction-stats-list.vue';

export default {
  inject: ['$system'],
  components: {
    RemediationInstructionStatsList,
  },
  mixins: [
    authMixin,
    localQueryMixin,
    entitiesRemediationInstructionStatsMixin,
    entitiesRemediationInstructionsMixin,
  ],
  data() {
    return {
      query: {
        interval: {
          from: QUICK_RANGES.last7Days.start,
          to: QUICK_RANGES.last7Days.stop,
        },
      },
    };
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    showRateInstructionModal(instruction) {
      this.$modals.show({
        name: MODALS.rate,
        config: {
          title: this.$t('modals.rateInstruction.title', { name: instruction.name }),
          text: this.$t('modals.rateInstruction.text'),
          action: async (data) => {
            await this.rateRemediationInstruction({ id: instruction._id, data });
            this.fetchList();
          },
        },
      });
    },

    fetchList() {
      const params = this.getQuery();
      params.with_flags = true;

      params.from = convertDateToTimestampByTimezone(dateParse(
        this.pagination.interval.from,
        DATETIME_INTERVAL_TYPES.start,
        DATETIME_FORMATS.datePicker,
      ), this.$system.timezone);
      params.to = convertDateToTimestampByTimezone(dateParse(
        this.pagination.interval.to,
        DATETIME_INTERVAL_TYPES.stop,
        DATETIME_FORMATS.datePicker,
      ), this.$system.timezone);

      this.fetchRemediationInstructionStatsList({ params });
    },
  },
};
</script>
