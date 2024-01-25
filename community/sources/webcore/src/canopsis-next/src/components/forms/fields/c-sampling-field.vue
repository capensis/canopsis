<template>
  <v-select
    v-field="value"
    :items="samplings"
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
  },
  computed: {
    samplings() {
      return Object.values(SAMPLINGS).map(value => ({
        value,
        text: this.$tc(`common.times.${value}`, 2),
      }));
    },
  },
};
</script>
