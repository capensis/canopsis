<template lang="pug">
  v-form
    v-switch(
      v-if="!isSimple",
      v-field="form.advancedMode",
      :label="$t('modals.eventFilterRule.advanced')",
      hide-details
    )
    v-text-field(
      v-field="form.field",
      v-validate="'required'",
      :label="$t('modals.eventFilterRule.field')",
      :error-messages="errors.collect('field')",
      name="field"
    )
    v-text-field(
      v-if="isSimple",
      v-field="form.value",
      :label="$t('modals.eventFilterRule.value')"
    )
    template(v-else)
      mixed-field(
        v-if="!form.advancedMode",
        v-field="form.value",
        :label="$t('modals.eventFilterRule.value')"
      )
      template(v-else)
        v-layout(align-center, justify-center)
          h2 {{ $t('modals.eventFilterRule.comparisonRules') }}
          v-btn(
            :disabled="!availableOperators.length > 0",
            icon,
            small,
            @click="addAdvancedRuleField"
          )
            v-icon add
        v-layout(v-for="(field, fieldIndex) in form.advancedRuleFields", :key="field.key", align-center)
          v-flex(xs3)
            v-select(
              v-field="form.advancedRuleFields[fieldIndex].key",
              v-validate="'required'",
              :items="getAvailableOperatorsForRule(field)",
              :error-messages="errors.collect(getFieldKeyName(field.key))",
              :name="getFieldKeyName(field.key)"
            )
          v-flex.pl-1(xs9)
            mixed-field(v-field="form.advancedRuleFields[fieldIndex].value")
          v-flex
            v-btn(small, icon, @click="deleteAdvancedRuleField(fieldIndex)")
              v-icon(color="error") delete
</template>

<script>
import formMixin from '@/mixins/form';

import MixedField from '@/components/forms/fields/mixed-field.vue';

export default {
  inject: ['$validator'],
  components: { MixedField },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    operators: {
      type: Array,
      default: () => [],
    },
    isSimple: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    availableOperators() {
      return this.operators.filter(operator => !this.form.advancedRuleFields.some(({ key }) => key === operator));
    },

    getAvailableOperatorsForRule() {
      return (rule) => {
        const rules = this.form.advancedRuleFields.filter(({ key }) => key !== rule.key);

        return this.operators.filter(operator => !rules.some(({ key }) => key === operator));
      };
    },

    getFieldKeyName() {
      return key => `fieldKey[${key}]`;
    },
  },
  methods: {
    addAdvancedRuleField() {
      const newAdvancedRuleFields = [
        ...this.form.advancedRuleFields,

        { key: this.availableOperators[0], value: '' },
      ];

      this.updateField('advancedRuleFields', newAdvancedRuleFields);
    },

    deleteAdvancedRuleField(field) {
      const newAdvancedRuleFields = this.form.advancedRuleFields.filter(({ key }) => key !== field.key);

      this.updateField('advancedRuleFields', newAdvancedRuleFields);
    },
  },
};
</script>

