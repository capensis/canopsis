<template lang="pug">
  c-advanced-data-table(
    :items="remediationInstructionChanges",
    :headers="headers",
    :loading="pending",
    :pagination.sync="pagination",
    :total-items="totalItems",
    advanced-pagination
  )
    template(slot="headerCell", slot-scope="props")
      span.c-table-header__text--multiline {{ props.header.text }}
    template(slot="modified_on", slot-scope="props")
      span {{ props.item.modified_on | date }}
    template(slot="execution_count", slot-scope="props")
      span {{ props.item.execution_count }}
    template(slot="avg_complete_time", slot-scope="props")
      span(v-if="props.item.execution_count") {{ props.item.avg_complete_time | duration }}
      span(v-else) {{ $t('remediationInstructionStats.notAvailable') }}
    template(slot="alarm_states", slot-scope="props")
      affect-alarm-states(v-if="props.item.execution_count", :alarm-states="props.item.alarm_states")
      template(v-else) -
    template(slot="ok_alarm_states", slot-scope="props")
      span.c-state-count-changes-chip.primary(v-if="props.item.execution_count") {{ props.item.ok_alarm_states }}
      template(v-else) -
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
          text: this.$t('remediationInstructionStats.modifiedOn'),
          value: 'modified_on',
          sortable: false,
        },
        {
          text: this.$t('remediationInstructionStats.averageCompletionTime'),
          value: 'avg_complete_time',
          sortable: false,
        },
        {
          text: this.$t('remediationInstructionStats.executionCount'),
          value: 'execution_count',
          sortable: false,
        },
        {
          text: this.$t('remediationInstructionStats.alarmStates'),
          value: 'alarm_states',
          sortable: false,
        },
        {
          text: this.$t('remediationInstructionStats.okAlarmStates'),
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
