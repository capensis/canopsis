<template lang="pug">
  v-layout(column)
    v-layout(v-for="(rule, index) in rules", :key="rule.key", row, justify-space-between, align-center)
      v-flex
        pattern-rule-field(
          v-bind="getRuleProps(rule)",
          :rule="rule",
          :name="rule.key",
          :attributes="attributes",
          @input="updateRule(index, $event)"
        )
      c-action-btn(
        :tooltip="$t('pattern.removeRule')",
        :disabled="readonly || disabled",
        type="delete",
        @click="removeItemFromArray(index)"
      )
    v-layout(v-if="!readonly", row, align-center)
      v-btn.ml-0(
        :disabled="disabled",
        :color="hasRulesErrors ? 'error' : 'primary'",
        outline,
        @click="addFilterRule"
      ) {{ $t('pattern.addRule') }}
      span.error--text(v-show="hasRulesErrors") {{ $t('pattern.errors.existExcluded') }}
</template>

<script>
import { patternRuleToForm, convertValueByOperator, getOperatorsByRule } from '@/helpers/entities/pattern/form';

import { formArrayMixin } from '@/mixins/form';

import PatternRuleField from './pattern-rule-field.vue';

export default {
  inject: ['$validator'],
  components: { PatternRuleField },
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
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    rulesMap() {
      return this.attributes.reduce((acc, { value, options = {} }) => {
        acc[value] = options;

        return acc;
      }, {});
    },

    hasRulesErrors() {
      return this.errors.has(this.name);
    },
  },
  created() {
    this.attachMinValueRule();
  },
  beforeDestroy() {
    this.detachMinValueRule();
  },
  methods: {
    attachMinValueRule() {
      this.$validator.attach({
        name: this.name,
        rules: {
          is_not: 0,
        },
        getter: () => {
          const rules = this.rules.filter((rule) => {
            const { disabled } = this.getOptionsByRule(rule);

            return !disabled;
          });

          return rules.length;
        },
        context: () => this,
        vm: this,
      });
    },

    detachMinValueRule() {
      this.$validator.detach(this.name);
    },

    getOptionsByRule(rule) {
      return this.rulesMap[rule.attribute] ?? {};
    },

    getRuleProps(rule) {
      const { operators, type, disabled, ...props } = this.getOptionsByRule(rule);

      return {
        ...props,
        disabled: disabled || this.readonly || this.disabled,
        type,
        operators: operators ?? getOperatorsByRule(rule, type),
      };
    },

    getUpdatedRule(rule, newRule) {
      const { defaultValue, operators } = this.getRuleProps(newRule);

      const updatedRule = { ...newRule };

      if (updatedRule.attribute !== rule.attribute) {
        updatedRule.operator = '';
        updatedRule.field = '';
        updatedRule.dictionary = '';
        updatedRule.value = defaultValue;
      } else if (updatedRule.operator !== rule.operator) {
        updatedRule.value = convertValueByOperator(updatedRule.value, updatedRule.operator);
      }

      if (updatedRule.value !== rule.value && operators?.length === 1) {
        [updatedRule.operator] = operators;
      }

      if (!operators.includes(updatedRule.operator)) {
        updatedRule.operator = '';
      }

      return updatedRule;
    },

    updateRule(index, newRule) {
      const rule = this.rules[index];

      this.updateItemInArray(index, this.getUpdatedRule(rule, newRule));
    },

    addFilterRule() {
      const [firstAttribute] = this.attributes;

      this.addItemIntoArray(patternRuleToForm({ field: firstAttribute?.value }));
    },
  },
};
</script>
