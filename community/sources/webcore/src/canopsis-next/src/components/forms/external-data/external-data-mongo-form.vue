<template lang="pug">
  v-layout(column)
    v-text-field(
      v-field="form.collection",
      v-validate="'required'",
      :label="$t('externalData.fields.collection')",
      :name="collectionFieldName",
      :error-messages="errors.collect(collectionFieldName)",
      :disabled="disabled"
    )
    v-layout(row)
      v-flex(xs6)
        v-text-field(
          v-field="form.sort_by",
          :label="$t('externalData.fields.sortBy')",
          :name="sortByFieldName",
          :error-messages="errors.collect(sortByFieldName)",
          :disabled="disabled"
        )
      v-flex.ml-3(xs6)
        v-select(
          v-field="form.sort",
          :items="sortOrders",
          :label="$t('externalData.fields.sort')",
          :name="sortFieldName",
          :error-messages="errors.collect(sortFieldName)",
          :disabled="disabled",
          clearable
        )
    external-data-mongo-condition-form(
      v-for="(condition, index) in form.conditions",
      v-field="form.conditions[index]",
      :key="condition.key",
      :name="`${name}.${condition.key}`",
      :disabled-remove="hasOnlyOneCondition",
      :disabled="disabled",
      :variables="variables",
      @remove="removeCondition(index)"
    )
    v-flex(v-if="!disabled")
      v-btn.ml-0.mb-0(color="primary", outline, @click="addCondition") {{ $t('common.addMore') }}
</template>

<script>
import { SORT_ORDERS } from '@/constants';

import { externalDataItemConditionAttributeToForm } from '@/helpers/forms/shared/external-data';

import { formMixin } from '@/mixins/form';

import ExternalDataMongoConditionForm from './external-data-mongo-condition-form.vue';

export default {
  inject: ['$validator'],
  components: { ExternalDataMongoConditionForm },
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
    name: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    variables: {
      type: Array,
      default: () => ([]),
    },
  },
  computed: {
    sortOrders() {
      return Object.values(SORT_ORDERS).map(order => order.toLowerCase());
    },

    hasOnlyOneCondition() {
      return this.form.conditions.length === 1;
    },

    collectionFieldName() {
      return `${this.name}.collection`;
    },

    sortFieldName() {
      return `${this.name}.sort`;
    },

    sortByFieldName() {
      return `${this.name}.sort_by`;
    },
  },
  methods: {
    addCondition() {
      this.updateField('conditions', [
        ...this.form.conditions,

        externalDataItemConditionAttributeToForm(),
      ]);
    },

    removeCondition(index) {
      this.updateField(
        'conditions',
        this.form.conditions.filter((condition, conditionIndex) => index !== conditionIndex),
      );
    },
  },
};
</script>
