<template lang="pug">
  div
    slot(v-if="isPatternsEmpty", name="no-data")
      v-alert.ma-2(
      :value="true",
      type="info"
      ) {{ disabled ? $t('patternsList.noDataDisabled') : $t('patternsList.noData') }}
    v-layout(
    v-for="(pattern, index) in patterns",
    :key="`${$options.filters.json(pattern)}${index}`",
    row,
    wrap,
    align-center
    )
      v-flex(:class="disabled ? 'xs12' : 'xs11'")
        v-textarea(
        :value="pattern | json",
        rows="7",
        no-resize,
        readonly,
        disabled
        )
      v-flex.text-xs-center(v-if="!disabled", xs1)
        div
          v-btn(icon, @click="showEditPatternModal(index)")
            v-icon edit
        div
          v-btn(color="error", icon, @click="showRemovePatternModal(index)")
            v-icon delete
    v-btn(v-if="!disabled", color="primary", @click="showCreatePatternModal") {{ $t('common.add') }}
</template>

<script>
import { MODALS, EVENT_FILTER_RULE_OPERATORS } from '@/constants';

import modalMixin from '@/mixins/modal';
import formArrayMixin from '@/mixins/form/array';

export default {
  mixins: [modalMixin, formArrayMixin],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Array,
      default: () => [],
    },
    operators: {
      type: Array,
      default: () => EVENT_FILTER_RULE_OPERATORS,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isPatternsEmpty() {
      return !this.patterns || !this.patterns.length;
    },
  },
  methods: {
    showCreatePatternModal() {
      this.showModal({
        name: MODALS.createEventFilterRulePattern,
        config: {
          operators: this.operators,
          action: pattern => this.addItemIntoArray(pattern),
        },
      });
    },

    showEditPatternModal(index) {
      this.showModal({
        name: MODALS.createEventFilterRulePattern,
        config: {
          pattern: this.patterns[index],
          operators: this.operators,
          action: pattern => this.updateItemInArray(index, pattern),
        },
      });
    },

    showRemovePatternModal(index) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: () => this.removeItemFromArray(index),
        },
      });
    },
  },
};
</script>
