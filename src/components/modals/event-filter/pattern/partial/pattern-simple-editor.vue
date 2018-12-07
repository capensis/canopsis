<template lang="pug">
  div
    v-layout(v-for="rule in pattern", :key="rule.key", justify-space-between)
      v-flex(xs6)
        v-text-field(v-model="rule.key")
      v-flex(v-if="isSimpleRule(rule.value)", xs5)
        v-text-field(v-model="rule.value")
      v-flex(v-else, xs5)
        div(v-for="value in rule.value", :key="value.key")
          v-layout(justify-space-between)
            v-flex(xs2)
              v-select(:items="operators", v-model="value.key")
            v-flex(xs9)
              v-text-field(v-model="value.value", :type="value.key !== 'regex' ? 'number' : null")
    v-btn(@click="convertPatternAndSave", color="primary") Save changes
</template>

<script>
export default {
  props: {
    value: {
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
      form: {
        advancedMode: false,
        field: '',
        value: '',
        advancedRuleFields: [],
      },
      pattern: [],
    };
  },
  watch: {
    value() {
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
    convertPatternToForm() {
      this.pattern = [];
      Object.keys(this.value).forEach((key) => {
        if (this.isSimpleRule(this.value[key])) {
          this.pattern.push({ key, value: this.value[key] });
        } else {
          const rule = [];
          Object.keys(this.value[key]).forEach((field) => {
            rule.push({ key: field, value: this.value[key][field] });
          });

          this.pattern.push({ key, value: rule });
        }
      });
    },
    convertPatternAndSave() {
      const pattern = [];
      this.pattern.forEach((rule) => {
        if (this.isSimpleRule(rule.value)) {
          pattern[rule.key] = rule.value;
        } else {
          const ruleValue = {};
          rule.value.map(value => ruleValue[value.key] = value.value);
          pattern[rule.key] = ruleValue;
        }
      });

      this.$emit('input', pattern);
    },
  },
};
</script>
