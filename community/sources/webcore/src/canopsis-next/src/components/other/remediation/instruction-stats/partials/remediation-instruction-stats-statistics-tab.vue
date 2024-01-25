<template>
  <c-advanced-data-table
    :items="remediationInstructionChanges"
    :headers="headers"
    :loading="pending"
    :options.sync="options"
    :total-items="totalItems"
    advanced-pagination
  >
    <template #headerCell="{ header }">
      <span class="c-table-header__text--multiline">{{ header.text }}</span>
    </template>
    <template #modified_on="{ item }">
      <span>{{ item.modified_on | date }}</span>
    </template>
    <template #execution_count="{ item }">
      <span>{{ item.execution_count }}</span>
    </template>
    <template #avg_complete_time="{ item }">
      <span v-if="item.execution_count">{{ item.avg_complete_time | duration }}</span>
      <span v-else>{{ $t('common.notAvailable') }}</span>
    </template>
    <template #alarm_states="{ item }">
      <affect-alarm-states
        v-if="item.execution_count"
        :alarm-states="item.alarm_states"
      />
      <template v-else>
        -
      </template>
    </template>
    <template #ok_alarm_states="{ item }">
      <c-state-count-changes-chip v-if="item.execution_count">
        {{ item.ok_alarm_states }}
      </c-state-count-changes-chip>
      <template v-else>
        -
      </template>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesRemediationInstructionStatsMixin } from '@/mixins/entities/remediation/instruction-stats';

import AffectAlarmStates from './affect-alarm-states.vue';

export default {
  components: { AffectAlarmStates },
  mixins: [localQueryMixin, entitiesRemediationInstructionStatsMixin],
  props: {
    remediationInstruction: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      remediationInstructionChanges: [],
      totalItems: 0,
      pending: false,
    };
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('remediation.instructionStat.modifiedOn'),
          value: 'modified_on',
          sortable: false,
        },
        {
          text: this.$t('remediation.instructionStat.averageCompletionTime'),
          value: 'avg_complete_time',
          sortable: false,
        },
        {
          text: this.$t('remediation.instructionStat.executionCount'),
          value: 'execution_count',
          sortable: false,
        },
        {
          text: this.$t('remediation.instructionStat.alarmStates'),
          value: 'alarm_states',
          sortable: false,
        },
        {
          text: this.$t('remediation.instructionStat.okAlarmStates'),
          value: 'ok_alarm_states',
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
        data: remediationInstructionChanges,
        meta,
      } = await this.fetchRemediationInstructionStatsChangesListWithoutStore({
        id: this.remediationInstruction._id,
        params: this.getQuery(),
      });

      this.remediationInstructionChanges = remediationInstructionChanges;
      this.totalItems = meta.total_count;
      this.pending = false;
    },
  },
};
</script>
