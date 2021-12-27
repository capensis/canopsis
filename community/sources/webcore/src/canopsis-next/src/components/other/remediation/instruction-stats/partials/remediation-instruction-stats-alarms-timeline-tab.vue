<template lang="pug">
  c-advanced-data-table(
    :items="remediationInstructionExecutions",
    :headers="headers",
    :loading="pending",
    :pagination.sync="pagination",
    :total-items="totalItems",
    advanced-pagination
  )
    template(slot="executed_on", slot-scope="props")
      span.c-nowrap {{ props.item.executed_on | date }}
    template(slot="alarm._id", slot-scope="props")
      span {{ props.item.alarm | get('v.display_name') }}
    template(slot="timeline", slot-scope="props")
      span.grey--text.text--darken-2(v-if="!props.item.alarm")
        | {{ $t('remediationInstructionStats.instructionChanged') }}
      alarm-horizontal-time-line.my-2(v-else, :alarm="props.item.alarm")
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
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      const {
        data: remediationInstructionExecutions,
        meta,
      } = await this.fetchRemediationInstructionStatsExecutionsListWithoutStore({
        id: this.remediationInstruction._id,
        params: this.getQuery(),
      });

      this.remediationInstructionExecutions = remediationInstructionExecutions;
      this.totalItems = meta.total_count;
      this.pending = false;
    },
  },
};
</script>
