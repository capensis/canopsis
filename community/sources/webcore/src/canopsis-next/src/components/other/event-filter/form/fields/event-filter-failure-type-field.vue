<template>
  <v-select
    v-field="value"
    :label="$t('eventFilter.filterByType')"
    :items="types"
    :error-messages="errors.collect(name)"
    :name="name"
    item-text="text"
    item-value="value"
    clearable
  />
</template>

<script>
import { EVENT_FILTER_FAILURE_TYPES } from '@/constants';

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
    name: {
      type: String,
      default: 'type',
    },
  },
  computed: {
    types() {
      return Object.values(EVENT_FILTER_FAILURE_TYPES).map(value => ({
        value,
        text: this.$t(`eventFilter.failureTypes.${value}`),
      }));
    },
  },
};
</script>
