<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.eventFilterRule.addAField') }}
    v-card-text
      v-form
        v-switch(:label="$t('modals.eventFilterRule.advanced')", v-model="form.advancedMode", hide-details)
        v-text-field(
        v-model="form.field",
        :label="$t('common.field')",
        name="field",
        v-validate="'required'",
        :error-messages="errors.collect('field')"
        )
        v-text-field(
        v-if="!form.advancedMode"
        v-model="form.value",
        :label="$t('common.value')",
        name="value",
        v-validate="'required'",
        :error-messages="errors.collect('value')"
        )
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
    v-layout.pa-2(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import modalInnerMixin from '@/mixins/modal/inner';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      pattern: {},
      form: {
        advancedMode: false,
        field: '',
        value: '',
        advancedRuleFields: [],
      },
    };
  },
  mounted() {
    if (this.config) {
      const {
        operators,
        ruleKey = '',
        ruleValue = '',
        isSimpleRule = true,
      } = { ...this.config };
      this.operators = operators;
      this.form.advancedMode = !isSimpleRule;
      this.form.field = ruleKey;

      if (isSimpleRule) {
        this.form.value = ruleValue;
      } else {
        this.form.advancedRuleFields = Object.keys(ruleValue).map(key => ({ key, value: ruleValue[key] }));
      }
    }
  },
  methods: {
    addAdvancedRuleField() {
      this.form.advancedRuleFields.push({ key: '<', value: '' });
    },
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          let newRule = {};

          if (!this.form.advancedMode) {
            newRule = { field: this.form.field, value: this.form.value };
          } else {
            const value = this.form.advancedRuleFields.reduce((acc, field) => {
              acc[field.key] = field.value;
              return acc;
            }, {});
            newRule = { field: this.form.field, value };
          }

          await this.config.action(newRule);
        }

        this.hideModal();
      }
    },
  },
};
</script>

