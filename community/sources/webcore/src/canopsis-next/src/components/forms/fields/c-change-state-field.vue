<template>
  <div class="mt-3">
    <v-layout>
      <state-criticity-field
        v-field="value.state"
        :state-values="availableStateValues"
      />
    </v-layout>
    <v-layout class="mt-4">
      <v-textarea
        v-field="value.output"
        v-validate="'required'"
        :label="label || $t('common.note')"
        :error-messages="errors.collect(outputFieldName)"
        :name="outputFieldName"
      />
    </v-layout>
  </div>
</template>

<script>
import { omit } from 'lodash';

import { ALARM_STATES } from '@/constants';

import { entitiesInfoMixin } from '@/mixins/entities/info';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

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
      return this.allowChangeSeverityToInfo ? ALARM_STATES : omit(ALARM_STATES, ['ok']);
    },

    outputFieldName() {
      return `${this.name}.output`;
    },
  },
};
</script>
