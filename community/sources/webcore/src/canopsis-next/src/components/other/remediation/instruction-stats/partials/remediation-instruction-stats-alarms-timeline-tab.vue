<template lang="pug">
  c-advanced-data-table(
    :items="remediationInstructionExecutions",
    :headers="headers",
    :loading="pending",
    :pagination.sync="pagination",
    :total-items="totalItems",
    advanced-pagination
  )
    template(#executed_on="{ item }")
      span.c-nowrap {{ item.executed_on | date }}
    template(#alarm._id="{ item }")
      span {{ item.alarm | get('v.display_name') }}
    template(#timeline="{ item }")
      span.grey--text.text--darken-2(v-if="!item.alarm") {{ $t('remediationInstructionStats.instructionChanged') }}
      alarm-horizontal-time-line.my-2(v-else, :alarm="item.alarm")
</template>

<script>
import { entitiesRemediationInstructionStatsMixin } from '@/mixins/entities/remediation/instruction-stats';
import { localQueryMixin } from '@/mixins/query-local/query';

import AlarmHorizontalTimeLine from '@/components/widgets/alarm/time-line/horizontal-time-line.vue';

export default {
  components: { AlarmHorizontalTimeLine },
  mixins: [localQueryMixin, entitiesRemediationInstructionStatsMixin],
  props: {
    remediationInstruction: {
      type: Object,
      default: () => ({}),
    },
    interval: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      remediationInstructionExecutions: [],
      totalItems: 0,
      pending: false,
    };
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('remediationInstructionStats.executedOn'),
          value: 'executed_on',
          sortable: false,
        },
        {
          text: this.$t('remediationInstructionStats.alarmId'),
          value: 'alarm.v.display_name',
          sortable: false,
        },
        {
          value: 'timeline',
          sortable: false,
        },
      ];
    },
  },
  watch: {
    interval: 'fetchList',
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      const params = this.getQuery();

      params.from = this.interval.from;
      params.to = this.interval.to;

      const {
        data: remediationInstructionExecutions,
        meta,
      } = await this.fetchRemediationInstructionStatsExecutionsListWithoutStore({
        params,
        id: this.remediationInstruction._id,
      });

      this.remediationInstructionExecutions = remediationInstructionExecutions;
      this.totalItems = meta.total_count;
      this.pending = false;
    },
  },
};
</script>
