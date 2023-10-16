<template>
  <v-select
    v-field="value"
    v-validate="rules"
    :label="label || $t('common.states')"
    :items="availableStateTypes"
    :disabled="disabled"
    :name="name"
    :error-messages="errors.collect(name)"
    hide-details
  />
</template>

<script>
import { ENTITIES_STATES } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Number,
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'state',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    availableStateTypes() {
      return Object.values(ENTITIES_STATES)
        .map(value => ({ value, text: this.$t(`common.stateTypes.${value}`) }));
    },
  },
};
</script>
