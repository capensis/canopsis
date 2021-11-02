<template lang="pug">
  v-layout
    v-flex(xs12)
      c-duration-field(
        v-field="timebased.time_interval",
        :label="$t('metaAlarmRule.timeInterval')",
        :units="availableUnits",
        required
      )
</template>

<script>
import { PERIODIC_REFRESH_UNITS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'timebased',
    event: 'input',
  },
  props: {
    timebased: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    availableUnits() {
      return Object.values(PERIODIC_REFRESH_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.timebased.time_interval.value),
      }));
    },
  },
};
</script>
