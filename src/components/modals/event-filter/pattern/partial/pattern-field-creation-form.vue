<template lang="pug">
  div
    v-form(ref="form")
      v-switch(:label="$t('modals.eventFilterRule.advanced')", v-model="form.advancedMode", hide-details)
      v-text-field(v-model="form.field", :label="$t('common.field')")
      template(v-if="!form.advancedMode")
        v-text-field(v-model="form.value", :label="$t('common.value')")
      template(v-else)
        v-layout(align-center, justify-center)
          h2 {{ $t('modals.eventFilterRule.comparisonRules') }}
          v-btn.primary(@click="addAdvancedRuleField", icon, fab, small)
            v-icon add
        v-layout(v-for="field in form.advancedRuleFields", :key="field.key")
          v-flex(xs3)
            v-select(:items="operators", v-model="field.key")
          v-flex(xs9)
            v-text-field(v-model="field.value")
    v-btn.primary(@click="submit") {{ $t('common.save') }}
</template>

<script>
import cloneDeep from 'lodash/cloneDeep';

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
    submit() {
      const newPattern = cloneDeep(this.value);
      if (!this.form.advancedMode) {
        newPattern[this.form.field] = this.form.value;
      } else {
        Object.values(this.form.advancedRuleFields)
          .reduce((acc, rule) => {
            acc[rule.key] = rule.value;
            return newPattern[this.form.field] = acc;
          }, {});
      }
      this.$emit('input', newPattern);
      this.$refs.form.reset();
    },
  },
};
</script>

