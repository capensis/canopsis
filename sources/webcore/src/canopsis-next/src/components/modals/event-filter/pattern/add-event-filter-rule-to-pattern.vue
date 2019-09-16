<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.eventFilterRule.addAField') }}
    v-card-text
      v-form
        v-switch(
        v-if="!config.isSimple",
        v-model="form.advancedMode",
        :label="$t('modals.eventFilterRule.advanced')",
        hide-details
        )
        v-text-field(
        v-model="form.field",
        :label="$t('modals.eventFilterRule.field')",
        name="field",
        v-validate="'required'",
        :error-messages="errors.collect('field')"
        )
        v-text-field(
        v-if="config.isSimple",
        v-model="form.value",
        :label="$t('modals.eventFilterRule.value')"
        )
        template(v-else)
          mixed-field(
          v-if="!form.advancedMode",
          v-model="form.value",
          :label="$t('modals.eventFilterRule.value')"
          )
          template(v-else)
            v-layout(align-center, justify-center)
              h2 {{ $t('modals.eventFilterRule.comparisonRules') }}
              v-btn(
              @click="addAdvancedRuleField",
              :disabled="!availableOperators.length > 0",
              icon,
              small,
              )
                v-icon add
            v-layout(v-for="field in form.advancedRuleFields", :key="field.key", align-center)
              v-flex(xs3)
                v-select(
                :items="getAvailableOperatorsForRule(field)",
                v-model="field.key",
                name="fieldKey",
                v-validate="'required'",
                :error-messages="errors.collect('fieldKey')"
                )
              v-flex.pl-1(xs9)
                mixed-field(v-model="field.value")
              v-flex
                v-btn(@click="deleteAdvancedRuleField(field)", small, icon)
                  v-icon(color="error") delete
    v-layout.pa-2(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { isObject } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import MixedField from '@/components/forms/fields/mixed-field.vue';

export default {
  name: MODALS.addEventFilterRuleToPattern,
  $_veeValidate: {
    validator: 'new',
  },
  components: { MixedField },
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
  computed: {
    availableOperators() {
      return this.operators.filter(operator => !this.form.advancedRuleFields.find(({ key }) => key === operator));
    },

    getAvailableOperatorsForRule() {
      return (rule) => {
        const rules = this.form.advancedRuleFields.filter(({ key }) => key !== rule.key);

        return this.operators.filter(operator => !rules.find(({ key }) => key === operator));
      };
    },
  },
  mounted() {
    if (this.config) {
      const {
        operators,
        ruleKey = '',
        ruleValue = '',
      } = this.config;

      const isSimpleRule = !isObject(ruleValue);

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
      this.form.advancedRuleFields.push({ key: this.availableOperators[0], value: '' });
    },

    deleteAdvancedRuleField(field) {
      this.form.advancedRuleFields = this.form.advancedRuleFields.filter(({ key }) => key !== field.key);
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

