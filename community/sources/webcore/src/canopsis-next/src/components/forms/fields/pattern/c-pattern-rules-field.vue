<template lang="pug">
  v-layout(column)
    v-layout(v-for="(rule, index) in rules", :key="rule.key", row, justify-space-between, align-center)
      v-flex
        c-pattern-rule-field(
          v-bind="getRuleProps(rule)",
          :rule="rule",
          :name="rule.key",
          :attributes="attributes",
          :disabled="disabled",
          @input="updateRule(index, $event)"
        )
      c-action-btn(
        :tooltip="$t('pattern.removeRule')",
        :disabled="disabled",
        type="delete",
        color="black",
        @click="removeItemFromArray(index)"
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
import { patternRuleToForm } from '@/helpers/forms/pattern';

import { convertValueByOperator, getOperatorsByRule } from '@/helpers/pattern';

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
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'rules',
    },
  },
  computed: {
    rulesMap() {
      return this.attributes.reduce((acc, { value, options = {} }) => {
        acc[value] = options;

        return acc;
      }, {});
    },
  },
  methods: {
    getOptionsByRule(rule) {
      return this.rulesMap[rule.attribute] ?? {};
    },

    getRuleProps(rule) {
      const { operators, type, ...props } = this.getOptionsByRule(rule);

      return {
        ...props,
        type,
        operators: operators ?? getOperatorsByRule(rule, type),
      };
    },

    getUpdatedRule(rule, newRule) {
      const { defaultValue } = this.getRuleProps(newRule);

      const updatedRule = { ...newRule };

      if (updatedRule.attribute !== rule.attribute) {
        updatedRule.operator = '';
        updatedRule.value = defaultValue ?? '';
      } else if (updatedRule.operator !== rule.operator) {
        updatedRule.value = convertValueByOperator(updatedRule.value, updatedRule.operator);
      }

      return updatedRule;
    },

    updateRule(index, newRule) {
      const rule = this.rules[index];

      this.updateItemInArray(index, this.getUpdatedRule(rule, newRule));
    },

    addFilterRule() {
      const [firstAttribute] = this.attributes;

      this.addItemIntoArray(patternRuleToForm({ attribute: firstAttribute?.value }));
    },
  },
};
</script>
