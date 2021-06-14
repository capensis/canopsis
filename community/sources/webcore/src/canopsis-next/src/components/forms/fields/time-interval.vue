<template lang="pug">
  v-layout.mb-3.time-interval(align-end)
    v-flex(xs6)
      v-text-field(
        v-field.number="interval.interval",
        v-validate="'required|numeric|min_value:1'",
        :label="intervalLabel || $t('common.interval')",
        :error-messages="errors.collect('interval')",
        :min="1",
        name="interval",
        type="number"
      )
    v-flex(xs6)
      v-select(
        v-field="interval.unit",
        v-validate="'required'",
        :label="unitLabel || $t('common.unit')",
        :error-messages="errors.collect('unit')",
        :items="availableUnits",
        name="unit"
      )
</template>

<script>
import { PERIODIC_REFRESH_UNITS, DEFAULT_TIME_INTERVAL } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'interval',
    event: 'input',
  },
  props: {
    interval: {
      type: Object,
      default: () => ({ ...DEFAULT_TIME_INTERVAL }),
    },
    intervalLabel: {
      type: String,
      required: false,
    },
    unitLabel: {
      type: String,
      required: false,
    },
    units: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    availableUnits() {
      if (this.units.length) {
        return this.units;
      }

      return Object.values(PERIODIC_REFRESH_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.interval.interval),
      }));
    },
  },
};
</script>

<style lang="scss" scoped>
.time-interval /deep/ .v-text-field__details {
  min-height: 14px;
}
</style>
