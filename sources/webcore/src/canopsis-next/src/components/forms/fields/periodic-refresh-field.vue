<template lang="pug">
  v-layout.mb-3(align-center)
    v-flex(xs5)
      enabled-field(
        v-field="periodicRefresh.enabled",
        :label="label",
        hide-details
      )
    v-flex(xs3)
      v-text-field(
        v-field="periodicRefresh.interval",
        v-validate="'required|numeric|min_value:0'",
        :error-messages="errors.collect('periodicRefreshInterval')",
        :disabled="!periodicRefresh.enabled",
        :min="0",
        name="periodicRefreshInterval",
        type="number",
        hide-details
      )
    v-flex(xs4)
      v-select(
        v-field="periodicRefresh.unit",
        v-validate="'required'",
        :items="availableUnits",
        :error-messages="errors.collect('periodicRefreshUnit')",
        :disabled="!periodicRefresh.enabled",
        name="periodicRefreshUnit",
        hide-details
      )
</template>

<script>
import { PERIODIC_REFRESH_UNITS, DEFAULT_PERIODIC_REFRESH } from '@/constants';

import EnabledField from '@/components/forms/fields/enabled-field.vue';

export default {
  components: { EnabledField },
  inject: ['$validator'],
  model: {
    prop: 'periodicRefresh',
    event: 'input',
  },
  props: {
    periodicRefresh: {
      type: Object,
      default: () => ({ ...DEFAULT_PERIODIC_REFRESH }),
    },
    label: {
      type: String,
      required: false,
    },
  },
  computed: {
    availableUnits() {
      return Object.values(PERIODIC_REFRESH_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, 2),
      }));
    },
  },
};
</script>
