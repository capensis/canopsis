<template>
  <v-select
    v-field="value"
    v-validate="'required'"
    :items="eventFilterTypes"
    :error-messages="errors.collect(name)"
    :label="label || $t('common.type')"
    :name="name"
  />
</template>

<script>
import { EVENT_FILTER_TYPES } from '@/constants';

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
      default: 'type',
    },
  },
  computed: {
    eventFilterTypes() {
      return Object.values(EVENT_FILTER_TYPES).map(type => ({
        value: type,
        text: this.$t(`eventFilter.types.${type}`),
      }));
    },
  },
};
</script>
