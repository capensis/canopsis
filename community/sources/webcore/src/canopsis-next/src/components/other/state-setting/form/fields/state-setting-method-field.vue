<template>
  <v-select
    v-field="value"
    v-validate="'required'"
    :items="methodTypes"
    :error-messages="errors.collect(name)"
    :label="label || $t('common.method')"
    :name="name"
  />
</template>

<script>
import { STATE_SETTING_METHODS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: STATE_SETTING_METHODS.worst,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'method',
    },
  },
  computed: {
    methodTypes() {
      return Object.values(STATE_SETTING_METHODS).map(method => ({
        value: method,
        text: this.$t(`stateSetting.methods.${method}`),
      }));
    },
  },
};
</script>
