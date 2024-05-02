<template>
  <v-select
    v-field="value"
    :items="availableSamplings"
    :disabled="disabled"
    :label="label || $t('common.sampling')"
    :name="name"
    hide-details
  >
    <template #selection="{ item }">
      <span class="text-capitalize">{{ item.text }}</span>
    </template>
    <template #item="{ item }">
      <span class="text-capitalize">{{ item.text }}</span>
    </template>
  </v-select>
</template>

<script>
import { SAMPLINGS } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'sampling',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    samplings: {
      type: Array,
      default: () => Object.values(SAMPLINGS),
    },
  },
  computed: {
    availableSamplings() {
      return this.samplings.map(value => ({
        value,
        text: this.$tc(`common.times.${value}`, 2),
      }));
    },
  },
};
</script>
