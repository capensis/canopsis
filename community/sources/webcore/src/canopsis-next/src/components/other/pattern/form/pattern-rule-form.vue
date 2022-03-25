<template lang="pug">
  div
    v-switch(
      v-if="!onlySimple",
      v-field="form.advancedMode",
      :label="$t('eventFilter.advanced')",
      color="primary",
      hide-details
    )
    v-text-field(
      v-model="form.field",
      v-validate="'required'",
      :label="$t('eventFilter.field')",
      :error-messages="errors.collect('field')",
      name="field"
    )
    v-text-field(
      v-if="onlySimple",
      v-model="form.value",
      :label="$t('eventFilter.value')"
    )
    template(v-else)
      c-mixed-field(
        v-if="!form.advancedMode",
        v-model="form.value",
        :label="$t('eventFilter.value')"
      )
      template(v-else)
        v-layout(align-center, justify-center)
          h2 {{ $t('eventFilter.comparisonRules') }}
          v-btn(
            :disabled="!availableOperators.length",
            icon,
            small,
            @click="addAdvancedField"
          )
            v-icon add
        pattern-rule-advanced-field-form(
          v-for="(field, index) in form.advancedFields",
          v-field="form.advancedFields[index]",
          :key="field.key",
          :operators="availableOperators",
          @delete="deleteAdvancedField(field)"
        )
</template>

<script>
import uid from '@/helpers/uid';

import { formMixin } from '@/mixins/form';

import PatternRuleAdvancedFieldForm from './pattern-rule-advanced-field-form.vue';

export default {
  inject: ['$validator'],
  components: { PatternRuleAdvancedFieldForm },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    operators: {
      type: Array,
      default: () => [],
    },
    onlySimple: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    availableOperators() {
      return this.operators.filter(
        operator => !this.form.advancedFields.some(({ operator: fieldOperator }) => fieldOperator === operator),
      );
    },
  },
  methods: {
    addAdvancedField() {
      const newField = { key: uid(), operator: this.availableOperators[0], value: '' };
      const newAdvancedFields = [
        ...this.form.advancedFields,

        newField,
      ];

      this.updateField('advancedFields', newAdvancedFields);
    },

    deleteAdvancedField(field) {
      const newAdvancedFields = this.form.advancedFields.filter(({ key }) => key !== field.key);

      this.updateField('advancedFields', newAdvancedFields);
    },
  },
};
</script>
