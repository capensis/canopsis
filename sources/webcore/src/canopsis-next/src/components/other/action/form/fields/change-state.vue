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
        :label="$t('modals.createAction.fields.output')"
      )
</template>

<script>
import { omit } from 'lodash';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

import { ENTITIES_STATES } from '@/constants';

export default {
  inject: ['$validator'],
  components: {
    StateCriticityField,
  },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  computed: {
    availableStateValues() {
      return omit(ENTITIES_STATES, ['ok']);
    },
  },
};
</script>
