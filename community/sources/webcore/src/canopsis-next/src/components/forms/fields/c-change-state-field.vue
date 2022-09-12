<template lang="pug">
  div.mt-3
    v-layout(row)
      state-criticity-field(
        v-field="value.state",
        :state-values="availableStateValues"
      )
    v-layout.mt-4(row)
      v-textarea(
        v-field="value.output",
        v-validate="'required'",
        :label="label || $t('common.note')",
        :error-messages="errors.collect(outputFieldName)",
        :name="outputFieldName"
      )
</template>

<script>
import { omit } from 'lodash';

import { entitiesInfoMixin } from '@/mixins/entities/info';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

import { ENTITIES_STATES } from '@/constants';

export default {
  inject: ['$validator'],
  components: {
    StateCriticityField,
  },
  mixins: [entitiesInfoMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'changeState',
    },
  },
  computed: {
    availableStateValues() {
      return this.allowChangeSeverityToInfo ? ENTITIES_STATES : omit(ENTITIES_STATES, ['ok']);
    },

    outputFieldName() {
      return `${this.name}.output`;
    },
  },
};
</script>
