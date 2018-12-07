<template lang="pug">
  div
    v-form
      v-switch(label="Advanced", v-model="form.advancedMode", hide-details)
      v-text-field(v-model="form.field", label="Field")
      template(v-if="!form.advancedMode")
        v-text-field(v-model="form.value", label="Value")
      template(v-else)
        v-layout(align-center, justify-center)
          h2 Comparison rules
          v-btn.primary(@click="addAdvancedRuleField", icon, fab, small)
            v-icon add
        v-layout(v-for="field in form.advancedRuleFields", :key="field.key")
          v-flex(xs3)
            v-select(:items="operators", v-model="field.key")
          v-flex(xs9)
            v-text-field(v-model="field.value")
    v-btn.primary(@click="addRule") Save
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
    };
  },
  methods: {
    addAdvancedRuleField() {
      this.form.advancedRuleFields.push({ key: '<', value: '' });
    },
    addRule() {
      const newPattern = { ...this.value };
      if (!this.form.advancedMode) {
        newPattern[this.form.field] = this.form.value;
      } else {
        const ruleValue = {};
        Object.values(this.form.advancedRuleFields)
          .forEach(rule => ruleValue[rule.key] = rule.value);
        newPattern[this.form.field] = ruleValue;
      }
      this.$emit('input', newPattern);
    },
  },
};
</script>

