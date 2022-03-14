<template lang="pug">
  v-layout(column)
    v-layout(v-for="(rule, index) in rules", :key="rule.key", row, justify-space-between, align-center)
      v-flex
        c-pattern-rule-field(
          v-field="rules[index]",
          v-bind="rulesMap[rule.attribute]",
          :attributes="attributes",
          :disabled="disabled"
        )
      c-action-btn(
        :tooltip="$t('pattern.removeRule')",
        :disabled="disabled",
        type="delete",
        color="black",
        @click="removeFilterRule(index)"
      )
    v-layout(row, align-center)
      v-btn.mx-0(
        :disabled="disabled",
        color="primary",
        outline,
        @click="addFilterRule"
      ) {{ $t('pattern.addRule') }}
</template>

<script>
import { filterRuleToForm } from '@/helpers/forms/filter';

import { formArrayMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formArrayMixin],
  model: {
    prop: 'rules',
    event: 'input',
  },
  props: {
    rules: {
      type: Array,
      required: true,
    },
    attributes: {
      type: Array,
      required: true,
    },
    rulesMap: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'rules',
    },
  },
  methods: {
    removeFilterRule(index) {
      if (this.rules.length !== 1) {
        this.removeItemFromArray(index);
      } else {
        this.$emit('remove');
      }
    },

    addFilterRule() {
      const [firstAttribute] = this.attributes;

      this.addItemIntoArray(filterRuleToForm({ attribute: firstAttribute?.value }));
    },
  },
};
</script>
