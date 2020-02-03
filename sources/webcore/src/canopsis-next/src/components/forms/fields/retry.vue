<template lang="pug">
  v-layout
    v-flex(xs4)
      v-text-field(
        v-field.number="value.count",
        v-validate="'required|numeric|min_value:0'",
        :label="$t('modals.createWebhook.fields.retryCount')",
        :error-messages="errors.collect('retryCount')",
        :min="0",
        name="retryCount",
        type="number"
      )
    v-flex(xs4)
      v-text-field(
        v-field.number="value.delay",
        v-validate="'required|numeric|min_value:0'",
        :label="$t('modals.createWebhook.fields.retryDelay')",
        :error-messages="errors.collect('retryDelay')",
        :min="0",
        name="retryDelay",
        type="number"
      )
    v-flex(xs4)
      v-select(
        v-field="value.unit",
        v-validate="'required'",
        :items="availableTypes",
        :error-messages="errors.collect('retryUnit')",
        name="retryUnit",
        hide-details
      )
        template(slot="selection", slot-scope="data")
          div.input-group__selections__comma {{ $tc(data.item.text, 2) }}
        template(slot="item", slot-scope="data")
          div.list__tile__title {{ $tc(data.item.text, 2) }}
</template>

<script>
import { PERIODIC_REFRESH_UNITS, DEFAULT_RETRY_FIELD } from '@/constants';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: Object,
      default: () => DEFAULT_RETRY_FIELD,
    },
  },
  computed: {
    availableTypes() {
      return Object.values(PERIODIC_REFRESH_UNITS);
    },
  },
};
</script>
