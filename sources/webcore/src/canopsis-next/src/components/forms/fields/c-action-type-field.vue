<template lang="pug">
  v-select(
    v-field="value",
    v-validate="'required'",
    :items="actionTypes",
    :error-messages="errors.collect(name)",
    :label="label || $t('common.type')",
    :name="name"
  )
</template>

<script>
import { SCENARIO_ACTION_TYPES } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'actionType',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    actionTypes() {
      return Object.values(SCENARIO_ACTION_TYPES).map(type => ({
        value: type,
        text: this.$t(`scenario.actions.${type}`),
      }));
    },
  },
};
</script>
