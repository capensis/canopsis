<template lang="pug">
  div.c-remediation-instruction-stats
    c-advanced-data-table(
      :headers="headers",
      :items="remediationInstructionStats",
      :loading="pending",
      :total-items="totalItems",
      :pagination="pagination",
      table-class="c-remediation-instruction-stats__table",
      expand,
      search,
      advanced-pagination,
      @update:pagination="$emit('update:pagination', $event)"
    )
      template(#toolbar="")
        v-layout(align-center)
          c-quick-date-interval-field(
            :interval="pagination.interval",
            :accumulated-before="accumulatedBefore",
            @input="updateInterval"
          )
      template(#headerCell="{ header }")
        span.c-table-header__text--multiline {{ header.text }}
      template(#type="{ item }") {{ $t(`remediation.instruction.types.${item.type}`) }}
      template(#last_executed_on="{ item }") {{ item.last_executed_on | date }}
      template(#last_modified="{ item }") {{ item.last_modified | date }}
      template(#avg_complete_time="{ item }") {{ item.avg_complete_time | duration }}
      template(#alarm_states="{ item }")
        affect-alarm-states(:alarm-states="item.alarm_states")
      template(#ok_alarm_states="{ item }")
        span.c-state-count-changes-chip.primary {{ item.ok_alarm_states }}
      template(#rating="{ item }")
        rating-field(:value="item.rating", readonly)
      template(#actions="{ item }")
        v-layout(row, justify-end)
          c-action-btn(
            v-if="item.rate_notify",
            :tooltip="$t('remediation.instructionStat.actions.needRate')",
            icon="notification_important",
            color="error",
            @click="$emit('rate', item)"
          )
          c-action-btn(
            :tooltip="$t('remediation.instructionStat.actions.rate')",
            :disabled="!item.ratable",
            icon="thumbs_up_down",
            @click="$emit('rate', item)"
          )
      template(#expand="{ item }")
        remediation-instruction-stats-list-expand-panel(
          :interval="interval",
          :remediation-instruction-stats-item="item"
        )
</template>

<script>
import { permissionsTechnicalRemediationInstructionMixin } from '@/mixins/permissions/technical/remediation-instruction';

import RatingField from '@/components/forms/fields/rating-field.vue';

import AffectAlarmStates from './partials/affect-alarm-states.vue';
import RemediationInstructionStatsListExpandPanel from './partials/remediation-instruction-stats-list-expand-panel.vue';

export default {
  components: {
    RemediationInstructionStatsListExpandPanel,
    RatingField,
    AffectAlarmStates,
  },
  mixins: [permissionsTechnicalRemediationInstructionMixin],
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
    accumulatedBefore: {
      type: Number,
      required: false,
    },
    pagination: {
      type: Object,
      required: true,
    },
    interval: {
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
          width: 150,
        },

        this.hasCreateAnyRemediationInstructionAccess && {
          text: this.$t('common.type'),
          value: 'type',
          width: 100,
        },

        {
          text: this.$t('remediation.instructionStat.lastExecutedOn'),
          value: 'last_executed_on',
          width: 180,
        },
        {
          text: this.$t('common.lastModifiedOn'),
          value: 'last_modified',
          width: 180,
        },
        {
          text: this.$t('remediation.instructionStat.averageCompletionTime'),
          value: 'avg_complete_time',
          sortable: false,
          width: 150,
        },
        {
          text: this.$t('remediation.instructionStat.executionCount'),
          value: 'execution_count',
          sortable: false,
          width: 150,
        },
        {
          text: this.$t('remediation.instructionStat.alarmStates'),
          value: 'alarm_states',
          sortable: false,
          width: 300,
        },
        {
          text: this.$t('remediation.instructionStat.okAlarmStates'),
          value: 'ok_alarm_states',
          sortable: false,
          width: 150,
        },
        {
          text: this.$tc('common.rating'),
          value: 'rating',
          sortable: false,
          width: 250,
        },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
          width: 100,
        },
      ].filter(Boolean);
    },
  },
  methods: {
    updateInterval(interval) {
      this.$emit('update:pagination', {
        ...this.pagination,
        interval,
      });
    },
  },
};
</script>

<style lang="scss">
.c-remediation-instruction-stats {
  &__table {
    table-layout: fixed !important;
  }
}
</style>
