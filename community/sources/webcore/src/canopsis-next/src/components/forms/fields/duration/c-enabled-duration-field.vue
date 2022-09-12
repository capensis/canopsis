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
        c-help-icon(v-if="helpText", slot="append", :text="helpText", max-width="300", top)
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
      div.v-messages.theme--light.error--text
        div.v-messages__wrapper
          div.v-messages__message(v-for="error in errors.collect(name)", :key="error") {{ error }}
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
  },
  computed: {
    enabledFieldName() {
      return `${this.name}.enabled`;
    },

    timeUnits() {
      return [
        AVAILABLE_TIME_UNITS.day,
        AVAILABLE_TIME_UNITS.week,
        AVAILABLE_TIME_UNITS.month,
        AVAILABLE_TIME_UNITS.year,
      ].map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.duration.value),
      }));
    },
  },
  created() {
    this.$validator.attach({
      name: this.name,
      context: () => this,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
  },
};
</script>