<template>
  <c-advanced-data-table
    :items="remediationInstructionExecutions"
    :headers="headers"
    :loading="pending"
    :options.sync="options"
    :total-items="totalItems"
    search
    advanced-pagination
  >
    <template #toolbar="">
      <v-layout align-center>
        <c-enabled-field
          v-model="showFailed"
          :label="$t('remediation.instructionStat.showFailedExecutions')"
          hide-details
        />
      </v-layout>
    </template>
    <template #executed_on="{ item }">
      <span class="c-nowrap">{{ (item.alarm ? item.executed_on : item.instruction_modified_on) | date }}</span>
    </template>
    <template #result="{ item }">
      <c-enabled
        v-if="item.alarm"
        :value="item.status === $constants.REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.completed"
      />
    </template>
    <template #duration="{ item }">
      <span>{{ item.duration | duration }}</span>
    </template>
    <template #resolved="{ item }">
      <span>{{ item.alarm | get('v.resolved') | date }}</span>
    </template>
    <template #timeout_after_execution="{ item }">
      <span>{{ item.timeout_after_execution | duration }}</span>
    </template>
    <template #timeline="{ item }">
      <span
        class="text--secondary"
        v-if="!item.alarm"
      >{{ $t('remediation.instructionStat.instructionChanged') }}</span>
      <alarm-horizontal-time-line
        class="my-2"
        v-else
        :alarm="item.alarm"
      />
    </template>
  </c-advanced-data-table>
</template>

<script>
import {
  prepareRemediationInstructionExecutionsForAlarmTimeline,
} from '@/helpers/entities/remediation/instruction-execution/list';

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
          text: this.$t('remediation.instructionStat.timeoutAfterExecution'),
          value: 'timeout_after_execution',
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
