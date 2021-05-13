<template lang="pug">
  div.instruction-list
    c-advanced-data-table.white(
      :headers="headers",
      :items="remediationInstructionStats",
      :loading="pending",
      :total-items="totalItems",
      :pagination="pagination",
      expand,
      search,
      advanced-pagination,
      @update:pagination="$emit('update:pagination', $event)"
    )
      template(slot="headerCell", slot-scope="props")
        span.pre-line.header-text {{ props.header.text }}
      template(slot="type", slot-scope="props") {{ $t(`remediationInstructions.types.${props.item.type}`) }}
      template(slot="last_executed_on", slot-scope="props") {{ props.item.last_executed_on | date('long', true, null) }}
      template(slot="last_modified", slot-scope="props") {{ props.item.last_modified | date('long', true, null) }}
      template(slot="avg_complete_time", slot-scope="props")
        | {{ props.item.avg_complete_time | duration }}
      template(slot="alarm_states", slot-scope="props")
        affect-alarm-states(:alarm-states="props.item.alarm_states")
      template(slot="ok_alarm_states", slot-scope="props")
        span.c-state-count-changes-chip.primary {{ props.item.ok_alarm_states }}
      template(slot="rating", slot-scope="props")
        rating-field(:value="props.item.rating", readonly)
      template(slot="actions", slot-scope="props")
        v-layout(row, justify-end)
          c-action-btn(
            v-if="props.item.rate_notify",
            :tooltip="$t('remediationInstructionStats.actions.needRate')",
            icon="notification_important",
            color="error",
            @click="$emit('approve', props.item)"
          )
          c-action-btn(
            v-if="props.item.execution_count > 0",
            :tooltip="$t('remediationInstructionStats.actions.rate')",
            icon="thumbs_up_down",
            @click="$emit('rate', props.item)"
          )
      template(slot="expand", slot-scope="props")
        remediation-instruction-stats-list-expand-panel(:remediation-instruction="props.item")
</template>

<script>
import RatingField from '@/components/forms/fields/rating-field.vue';

import AffectAlarmStates from './partials/affect-alarm-states.vue';
import RemediationInstructionStatsListExpandPanel from './partials/remediation-instruction-stats-list-expand-panel.vue';

export default {
  components: { RemediationInstructionStatsListExpandPanel, RatingField, AffectAlarmStates },
  props: {
    remediationInstructionStats: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pagination: {
      type: Object,
      required: true,
    },
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.type'),
          value: 'type',
        },
        {
          text: this.$t('remediationInstructionStats.table.lastExecutedOn'),
          value: 'last_executed_on',
        },
        {
          text: this.$t('remediationInstructionStats.table.lastModifiedOn'),
          value: 'last_modified',
        },
        {
          text: this.$t('remediationInstructionStats.table.averageCompletionTime'),
          value: 'avg_complete_time',
          sortable: false,
        },
        {
          text: this.$t('remediationInstructionStats.table.executionCount'),
          value: 'execution_count',
          sortable: false,
        },
        {
          text: this.$t('remediationInstructionStats.table.alarmStates'),
          value: 'alarm_states',
          sortable: false,
        },
        {
          text: this.$t('remediationInstructionStats.table.okAlarmStates'),
          value: 'ok_alarm_states',
          sortable: false,
        },
        {
          text: this.$t('remediationInstructionStats.table.rating'),
          value: 'rating',
          sortable: false,
        },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ];
    },
  },
};
</script>

<style lang="scss" scoped>
.header-text {
  display: inline-block;
  height: 100%;
  vertical-align: middle;
}

.instruction-list {
  /deep/ thead th {
    vertical-align: middle;
  }
}
</style>
