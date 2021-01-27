<template lang="pug">
  v-select(
    v-field="value",
    :label="$t('webhook.disableDuringPeriods')",
    :items="duringItems",
    multiple
  )
</template>

<script>
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      required: false,
    },
  },
  computed: {
    duringItems() {
      return [PBEHAVIOR_TYPE_TYPES.maintenance, PBEHAVIOR_TYPE_TYPES.pause];
    },
  },
};
</script>
