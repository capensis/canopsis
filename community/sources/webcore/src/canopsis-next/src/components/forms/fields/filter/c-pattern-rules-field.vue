<template lang="pug">
  v-layout(column)
    v-layout(v-for="(rule, index) in rules", :key="rule.key", row, justify-space-between, align-center)
      v-flex
        c-pattern-rule-field(v-field="rules[index]", v-bind="getFilterRuleProps(rule)")
      c-action-btn(
        :tooltip="$t('patterns.removeRule')",
        type="delete",
        color="black",
        @click="removeFilterRule(index)"
      )
    v-layout
      v-btn.mx-0(color="primary", outline, @click="addFilterRule") {{ $t('patterns.addRule') }}
</template>

<script>
import { filterRuleToForm } from '@/helpers/forms/filter';

import { formArrayMixin } from '@/mixins/form';

export default {
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

    getFilterRuleProps(rule) {
      const props = this.rulesMap[rule.attribute] ?? {};

      return {
        attributes: this.attributes,
        ...props,
      };
    },
  },
};
</script>
