<template lang="pug">
  div.mt-3
    v-layout(row)
      state-criticity-field(
        v-field="value.state",
        :stateValues="availableStateValues"
      )
    v-layout.mt-4(row)
      v-textarea(
        v-field="value.output",
        :label="label || $t('modals.createAction.fields.output')",
        v-validate="'required'",
        :error-messages="errors.collect('output')",
        name="output"
      )
</template>

<script>
import { omit } from 'lodash';

import entitiesInfoMixin from '@/mixins/entities/info';

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
      required: false,
    },
  },
  computed: {
    availableStateValues() {
      return this.allowChangeSeverityToInfo ? ENTITIES_STATES : omit(ENTITIES_STATES, ['ok']);
    },
  },
};
</script>
