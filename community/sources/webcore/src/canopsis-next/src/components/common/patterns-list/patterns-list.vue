<template lang="pug">
  div
    slot(v-if="isPatternsEmpty", name="no-data")
      v-alert(
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
        v-layout
          pattern-information.ma-3.ml-0 {{ $t('common.and') }}
          v-flex(xs12)
            v-textarea(
              :value="pattern | json",
              rows="7",
              no-resize,
              readonly,
              disabled
            )
        v-layout(v-if="index !== patterns.length - 1", justify-center)
          span.text-uppercase.operator-chip {{ $t('common.or') }}
      v-flex.text-xs-center(v-if="!disabled", xs1)
        div
          v-btn(icon, @click="showEditPatternModal(index)")
            v-icon edit
        div
          v-btn(color="error", icon, @click="showRemovePatternModal(index)")
            v-icon delete
    v-btn.mx-0(v-if="!disabled", color="primary", @click="showCreatePatternModal") {{ $t('common.add') }}
    v-layout(v-if="errors", row)
      v-alert(:value="errors.has(name)", type="error")
        span(v-for="error in errors.collect(name)", :key="error") {{ error }}
</template>

<script>
import { MODALS, EVENT_FILTER_RULE_OPERATORS } from '@/constants';

import formArrayMixin from '@/mixins/form/array';

import PatternInformation from '@/components/other/pattern/pattern-information.vue';

export default {
  $_veeValidate: {
    value() {
      return this.patterns;
    },

    name() {
      return this.name;
    },
  },
  inject: ['$validator'],
  components: {
    PatternInformation,
  },
  mixins: [formArrayMixin],
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
    name: {
      type: String,
      default: 'patterns',
    },
  },
  computed: {
    isPatternsEmpty() {
      return !this.patterns || !this.patterns.length;
    },
  },
  methods: {
    showCreatePatternModal() {
      this.$modals.show({
        name: MODALS.createPattern,
        config: {
          operators: this.operators,
          action: pattern => this.addItemIntoArray(pattern),
        },
      });
    },

    showEditPatternModal(index) {
      this.$modals.show({
        name: MODALS.createPattern,
        config: {
          pattern: this.patterns[index],
          operators: this.operators,
          action: pattern => this.updateItemInArray(index, pattern),
        },
      });
    },

    showRemovePatternModal(index) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.removeItemFromArray(index),
        },
      });
    },
  },
};
</script>
