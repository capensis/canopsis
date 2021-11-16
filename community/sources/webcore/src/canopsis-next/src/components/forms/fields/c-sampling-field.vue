<template lang="pug">
  v-select(
    v-field="value",
    :items="samplings",
    :disabled="disabled",
    :label="label || $t('common.sampling')",
    :name="name",
    hide-details
  )
    template(#selection="{ item }")
      span.text-capitalize {{ item.text }}
    template(#item="{ item }")
      span.text-capitalize {{ item.text }}
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
      required: true,
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
