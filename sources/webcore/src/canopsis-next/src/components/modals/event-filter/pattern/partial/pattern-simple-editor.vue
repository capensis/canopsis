<template lang="pug">
  div
    v-layout(v-for="(rule, ruleKey) in convertedPattern", :key="rule.key", justify-space-between)
      v-flex(xs6)
        v-text-field(v-model="rule.key")
      v-flex(v-if="isSimpleRule(rule.value)", xs5)
        v-text-field(v-model="rule.value")
      v-flex(v-else, xs5)
        div(v-for="(value, subRuleKey) in rule.value", :key="value.key")
          v-layout(justify-space-between)
            v-flex(xs2)
              v-select(
              :items="operators",
              :value="value.key",
              @input="editAdvancedRuleOperator($event, value, ruleKey, subRuleKey)"
              )
            v-flex(xs9)
              v-text-field(
              :value="value.value",
              @input="editAdvancedRuleValue($event, value, ruleKey, subRuleKey)"
              :type="value.key !== 'regex' ? 'number' : null"
              )
    v-btn(@click="convertPatternAndSave", color="primary") {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

export default {
  model: {
    prop: 'pattern',
    event: 'input',
  },
  props: {
    pattern: {
      type: Object,
      required: true,
    },
    operators: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      convertedPattern: [],
    };
  },
  watch: {
    pattern() {
      this.convertPatternToForm();
    },
  },
  methods: {
    isSimpleRule(rule) {
      if (typeof rule === 'string') {
        return true;
      }

      return false;
    },

    editAdvancedRuleOperator(operator, rule, ruleKey, subRuleKey) {
      const newPattern = [...this.convertedPattern];
      newPattern[ruleKey].value[subRuleKey] = { ...rule, key: operator };
      this.convertedPattern = newPattern;
    },

    editAdvancedRuleValue(value, rule, ruleKey, subRuleKey) {
      const newPattern = [...this.convertedPattern];
      newPattern[ruleKey].value[subRuleKey] = { ...rule, value };
      this.convertedPattern = newPattern;
    },

    convertPatternToForm() {
      const pattern = cloneDeep(this.pattern);
      this.convertedPattern = Object.keys(pattern).reduce((acc, key) => {
        if (this.isSimpleRule(pattern[key])) {
          acc.push({ key, value: cloneDeep(pattern[key]) });
        } else {
          acc.push({ key, value: this.convertAdvancedRuleToForm(pattern[key]) });
        }
        return acc;
      }, []);
    },

    convertAdvancedRuleToForm(rule) {
      return Object.keys(rule).reduce((acc, key) => {
        acc.push({ key, value: rule[key] });
        return acc;
      }, []);
    },

    convertAdvancedRuleToObject(rule) {
      return rule.value.reduce((acc, key) => {
        acc[key.key] = key.value;
        return acc;
      }, {});
    },

    convertPatternAndSave() {
      const newPattern = this.convertedPattern.reduce((acc, rule) => {
        if (this.isSimpleRule(rule.value)) {
          acc[rule.key] = rule.value;
        } else {
          acc[rule.key] = this.convertAdvancedRuleToObject(rule);
        }

        return acc;
      }, {});

      this.$emit('input', newPattern);
    },
  },
};
</script>
