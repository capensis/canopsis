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
      template(slot="toolbar", slot-scope="props")
        v-layout(align-center)
          c-quick-date-interval-field(
            :interval="pagination.interval",
            :accumulated-before="accumulatedBefore",
            @input="updateInterval"
          )
      template(slot="headerCell", slot-scope="props")
        span.c-table-header__text--multiline {{ props.header.text }}
      template(slot="type", slot-scope="props")
        | {{ $t(`remediationInstructions.types.${props.item.type}`) }}
      template(slot="last_executed_on", slot-scope="props")
        | {{ props.item.last_executed_on | date }}
      template(slot="last_modified", slot-scope="props")
        | {{ props.item.last_modified | date }}
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
            @click="$emit('rate', props.item)"
          )
          c-action-btn(
            :tooltip="$t('remediationInstructionStats.actions.rate')",
            :disabled="!props.item.ratable",
            icon="thumbs_up_down",
            @click="$emit('rate', props.item)"
          )
      template(slot="expand", slot-scope="props")
        remediation-instruction-stats-list-expand-panel(:remediation-instruction-stats-id="props.item._id")
</template>

<script>
import {
  permissionsTechnicalRemediationInstructionMixin,
} from '@/mixins/permissions/technical/remediation-instruction';

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
          text: this.$t('remediationInstructionStats.lastExecutedOn'),
          value: 'last_executed_on',
          width: 180,
        },
        {
          text: this.$t('common.lastModifiedOn'),
          value: 'last_modified',
          width: 180,
        },
        {
          text: this.$t('remediationInstructionStats.averageCompletionTime'),
          value: 'avg_complete_time',
          sortable: false,
          width: 150,
        },
        {
          text: this.$t('remediationInstructionStats.executionCount'),
          value: 'execution_count',
          sortable: false,
          width: 150,
        },
        {
          text: this.$t('remediationInstructionStats.alarmStates'),
          value: 'alarm_states',
          sortable: false,
          width: 300,
        },
        {
          text: this.$t('remediationInstructionStats.okAlarmStates'),
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
