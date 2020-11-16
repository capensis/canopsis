<template lang="pug">
  div.instruction-list
    advanced-data-table.white(
      :headers="headers",
      :items="remediationInstructions",
      :loading="pending",
      :total-items="totalItems",
      :pagination="pagination",
      :is-disabled-item="isDisabledInstruction",
      select-all,
      expand,
      search,
      advanced-pagination,
      @update:pagination="$emit('update:pagination', $event)"
    )
      template(slot="toolbar", slot-scope="props")
        v-flex(v-show="hasDeleteAnyRemediationInstructionAccess && props.selected.length", xs4)
          v-btn(@click="$emit('remove-selected', props.selected)", icon)
            v-icon delete
      template(slot="headerCell", slot-scope="props")
        span.pre-line.header-text {{ props.header.text }}
      template(slot="enabled", slot-scope="props")
        enabled-column(:value="props.item.enabled")
      template(slot="rating", slot-scope="props")
        rating-field(:value="props.item.rating", readonly)
      template(slot="last_modified", slot-scope="props")
        | {{ props.item.last_modified | date('long', true, null) }}
      template(slot="avg_complete_time", slot-scope="props")
        span(v-if="props.item.avg_complete_time")
          | {{ props.item.avg_complete_time }}
          | {{ $tc('common.times.minute', props.item.avg_complete_time) }}
      template(slot="last_executed_on", slot-scope="props")
        | {{ props.item.last_executed_on | date('long', true, null) }}
      template(slot="actions", slot-scope="props")
        v-layout
          v-btn.mx-0(
            v-if="hasUpdateAnyRemediationInstructionAccess",
            icon,
            small,
            @click.stop="$emit('edit', props.item)"
          )
            v-icon edit
          v-btn.mx-0(
            v-if="hasUpdateAnyRemediationInstructionAccess",
            icon,
            small,
            @click.stop="$emit('assign-patterns', props.item)"
          )
            v-icon assignment
          v-tooltip(bottom, :disabled="!props.disabled")
            v-btn.mx-0(
              slot="activator",
              v-if="hasDeleteAnyRemediationInstructionAccess",
              :disabled="props.disabled",
              icon,
              small,
              @click.stop="$emit('remove', props.item)"
            )
              v-icon(color="error") delete
            span {{ $t('remediationInstructions.usingInstruction') }}
      template(slot="expand", slot-scope="props")
        remediation-instructions-list-expand-panel(:remediationInstruction="props.item")
</template>

<script>
import rightsTechnicalRemediationInstructionMixin from '@/mixins/rights/technical/remediation-instruction';

import EnabledColumn from '@/components/tables/enabled-column.vue';
import RatingField from '@/components/forms/fields/rating-field.vue';

import RemediationInstructionsListExpandPanel from './partials/remediation-instructions-list-expand-panel.vue';

export default {
  components: {
    RatingField,
    EnabledColumn,
    RemediationInstructionsListExpandPanel,
  },
  mixins: [rightsTechnicalRemediationInstructionMixin],
  props: {
    remediationInstructions: {
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
          text: this.$t('common.author'),
          value: 'author',
        },
        {
          text: this.$t('common.enabled'),
          value: 'enabled',
        },
        {
          text: this.$t('remediationInstructions.table.rating'),
          value: 'rating',
        },
        {
          text: this.$t('remediationInstructions.table.lastModifiedOn'),
          value: 'last_modified',
        },
        {
          text: this.$t('remediationInstructions.table.averageTimeCompletion'),
          value: 'avg_complete_time',
        },
        {
          text: this.$t('remediationInstructions.table.monthExecutions'),
          value: 'month_executions',
        },
        {
          text: this.$t('remediationInstructions.table.lastExecutedBy'),
          value: 'last_executed_by.username',
        },
        {
          text: this.$t('remediationInstructions.table.lastExecutedOn'),
          value: 'last_executed_on',
        },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ];
    },
  },
  methods: {
    isDisabledInstruction({ deletable }) {
      return !deletable;
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
