<template lang="pug">
  v-layout(row, wrap)
    v-flex(xs5)
      v-checkbox(
        v-validate,
        v-field="duration.enabled",
        :label="label",
        :error-messages="errors.collect(enabledFieldName)",
        :name="enabledFieldName",
        color="primary"
      )
        template(#append="")
          c-help-icon(v-if="helpText", :text="helpText", color="info", max-width="300", top)
    v-flex(xs4)
      c-duration-field(
        v-field="duration",
        :units-label="$t('common.unit')",
        :disabled="!duration.enabled",
        :required="duration.enabled",
        :units="timeUnits",
        :name="name"
      )
    v-flex(xs9)
      v-messages(:value="errors.collect(name)", color="error")
</template>

<script>
import { AVAILABLE_TIME_UNITS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'duration',
    event: 'input',
  },
  props: {
    duration: {
      type: Object,
      required: true,
    },
    label: {
      type: String,
      required: true,
    },
    helpText: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      required: false,
    },
    units: {
      type: Array,
      required: false,
    },
  },
  computed: {
    enabledFieldName() {
      return `${this.name}.enabled`;
    },

    timeUnits() {
      const units = this.units || [
        AVAILABLE_TIME_UNITS.day,
        AVAILABLE_TIME_UNITS.week,
        AVAILABLE_TIME_UNITS.month,
        AVAILABLE_TIME_UNITS.year,
      ];

      return units.map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.duration.value),
      }));
    },
  },
  created() {
    this.$validator.attach({
      name: this.name,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
  },
};
</script>
