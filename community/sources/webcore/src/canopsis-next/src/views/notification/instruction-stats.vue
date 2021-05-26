<template lang="pug">
  div
    c-the-page-header {{ $t('common.instructionRating') }}
    remediation-instruction-stats-list(
      :remediation-instruction-stats="remediationInstructionStats",
      :pending="remediationInstructionStatsPending",
      :pagination.sync="pagination",
      :total-items="remediationInstructionStatsMeta.total_count",
      @rate="showInstructionRateModal"
    )
    c-fab-btn(@refresh="fetchList")
</template>

<script>
import moment from 'moment';

import { DATETIME_FORMATS, DATETIME_INTERVAL_TYPES } from '@/constants';

import { dateParse } from '@/helpers/date/date-intervals';

import { authMixin } from '@/mixins/auth';
import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesRemediationInstructionStatsMixin } from '@/mixins/entities/remediation/instruction-stats';

import RemediationInstructionStatsList from '@/components/other/remediation/instruction-stats/remediation-instruction-stats-list.vue';

export default {
  components: {
    RemediationInstructionStatsList,
  },
  mixins: [
    authMixin,
    localQueryMixin,
    entitiesRemediationInstructionStatsMixin,
  ],
  data() {
    return {
      query: {
        interval: {
          from: moment().subtract(1, 'week').format(DATETIME_FORMATS.datePicker),
          to: moment().format(DATETIME_FORMATS.datePicker),
        },
      },
    };
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    showInstructionRateModal() {},

    fetchList() {
      const params = this.getQuery();
      params.with_flags = true;

      params.from = dateParse(
        this.pagination.interval.from,
        DATETIME_INTERVAL_TYPES.start,
        DATETIME_FORMATS.datePicker,
      ).unix();
      params.to = dateParse(
        this.pagination.interval.to,
        DATETIME_INTERVAL_TYPES.stop,
        DATETIME_FORMATS.datePicker,
      ).unix();

      this.fetchRemediationInstructionStatsList({ params });
    },
  },
};
</script>
