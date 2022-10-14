<template lang="pug">
  v-layout(column)
    v-layout(row)
      v-text-field(
        v-field="form.collection",
        v-validate="'required'",
        :label="$t('eventFilter.collection')",
        :name="collectionFieldName",
        :error-messages="errors.collect(collectionFieldName)",
        :disabled="disabled"
      )
    event-filter-enrichment-external-data-mongo-condition-form(
      v-for="(condition, index) in form.conditions",
      v-field="form.conditions[index]",
      :key="condition.key",
      :name="`${name}.${condition.key}`",
      :disabled-remove="hasOnlyOneCondition",
      :disabled="disabled",
      @remove="removeCondition(index)"
    )
    v-flex(v-if="!disabled")
      v-btn.ml-0.mb-0(color="primary", outline, @click="addCondition") {{ $t('common.addMore') }}
</template>

<script>
import { eventFilterExternalDataConditionItemToForm } from '@/helpers/forms/event-filter';

import { formMixin } from '@/mixins/form';

import EventFilterEnrichmentExternalDataMongoConditionForm
  from './event-filter-enrichment-external-data-mongo-condition-form.vue';

export default {
  inject: ['$validator'],
  components: { EventFilterEnrichmentExternalDataMongoConditionForm },
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
  },
  computed: {
    hasOnlyOneCondition() {
      return this.form.conditions.length === 1;
    },

    collectionFieldName() {
      return `${this.name}.collection`;
    },
  },
  methods: {
    addCondition() {
      this.updateField('conditions', [
        ...this.form.conditions,

        eventFilterExternalDataConditionItemToForm(),
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
