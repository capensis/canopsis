<template lang="pug">
  c-advanced-data-table(
    :items="remediationInstructionExecutions",
    :headers="headers",
    :loading="pending",
    :pagination.sync="pagination",
    :total-items="totalItems",
    search,
    advanced-pagination
  )
    template(#toolbar="")
      v-layout(align-center)
        c-enabled-field(
          v-model="showFailed",
          :label="$t('remediation.instructionStat.showFailedExecutions')",
          hide-details
        )
    template(#executed_on="{ item }")
      span.c-nowrap {{ (item.alarm ? item.executed_on : item.instruction_modified_on) | date }}
    template(#result="{ item }")
      c-enabled(
        v-if="item.alarm",
        :value="item.status === $constants.REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.completed"
      )
    template(#duration="{ item }")
      span {{ item.duration | duration }}
    template(#resolved="{ item }")
      span {{ item.alarm | get('v.resolved') | date }}
    template(#timeline="{ item }")
      span.grey--text.text--darken-2(v-if="!item.alarm") {{ $t('remediation.instructionStat.instructionChanged') }}
      alarm-horizontal-time-line.my-2(v-else, :alarm="item.alarm")
</template>

<script>
import {
  prepareRemediationInstructionExecutionsForAlarmTimeline,
} from '@/helpers/entities/remediation-instruction-execution';

import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesRemediationInstructionStatsMixin } from '@/mixins/entities/remediation/instruction-stats';

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
      showFailed: true,
      pending: false,
    };
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('remediation.instructionStat.executedAt'),
          value: 'executed_on',
          sortable: false,
        },
        {
          text: this.$t('common.alarmId'),
          value: 'alarm.v.display_name',
          sortable: false,
        },
        {
          text: this.$t('common.result'),
          value: 'result',
          sortable: false,
        },
        {
          text: this.$t('remediation.instructionStat.remediationDuration'),
          value: 'duration',
          sortable: false,
        },
        {
          text: this.$t('remediation.instructionStat.alarmResolvedDate'),
          value: 'resolved',
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
    showFailed: 'fetchList',
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
      params.show_failed = this.showFailed;

      const {
        data: remediationInstructionExecutions,
        meta,
      } = await this.fetchRemediationInstructionStatsExecutionsListWithoutStore({
        params,
        id: this.remediationInstruction._id,
      });

      this.remediationInstructionExecutions = prepareRemediationInstructionExecutionsForAlarmTimeline(
        remediationInstructionExecutions,
      );
      this.totalItems = meta.total_count;
      this.pending = false;
    },
  },
};
</script>
