<template lang="pug">
  div
    c-page-header
    v-card.ma-4.mt-0
      remediation-instruction-stats-list(
        :remediation-instruction-stats="remediationInstructionStats",
        :pending="remediationInstructionStatsPending",
        :pagination.sync="pagination",
        :total-items="remediationInstructionStatsMeta.total_count",
        :accumulated-before="remediationInstructionStatsMeta.accumulated_before",
        :interval="interval",
        @rate="showRateInstructionModal"
      )
    c-fab-btn(@refresh="fetchList")
</template>

<script>
import { MODALS, QUICK_RANGES } from '@/constants';

import {
  convertStartDateIntervalToTimestamp,
  convertStopDateIntervalToTimestamp,
} from '@/helpers/date/date-intervals';
import {
  convertDateToEndOfDayTimestamp,
  convertDateToStartOfDayTimestamp,
} from '@/helpers/date/date';

import { authMixin } from '@/mixins/auth';
import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesRemediationInstructionStatsMixin } from '@/mixins/entities/remediation/instruction-stats';
import { entitiesRemediationInstructionMixin } from '@/mixins/entities/remediation/instruction';

import RemediationInstructionStatsList from '@/components/other/remediation/instruction-stats/remediation-instruction-stats-list.vue';

export default {
  components: {
    RemediationInstructionStatsList,
  },
  mixins: [
    authMixin,
    localQueryMixin,
    entitiesRemediationInstructionStatsMixin,
    entitiesRemediationInstructionMixin,
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
  computed: {
    interval() {
      return {
        from: convertDateToStartOfDayTimestamp(convertStartDateIntervalToTimestamp(
          this.query.interval.from,
        )),
        to: convertDateToEndOfDayTimestamp(convertStopDateIntervalToTimestamp(
          this.query.interval.to,
        )),
      };
    },
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
      params.from = this.interval.from;
      params.to = this.interval.to;

      this.fetchRemediationInstructionStatsList({ params });
    },
  },
};
</script>
