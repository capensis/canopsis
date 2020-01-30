<template lang="pug">
  v-layout.mb-3(align-center)
    v-flex(xs5)
      v-switch(
        v-model="periodicRefresh.enabled",
        :label="label",
        hide-details
      )
    v-flex(xs3)
      v-text-field(
        v-model="periodicRefresh.interval",
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
        v-model="periodicRefresh.unit",
        v-validate="'required'",
        :items="availableTypes",
        :error-messages="errors.collect('periodicRefreshUnit')",
        :disabled="!periodicRefresh.enabled",
        name="periodicRefreshUnit",
        hide-details
      )
        template(slot="selection", slot-scope="data")
          div.input-group__selections__comma {{ $tc(data.item.text, 2) }}
        template(slot="item", slot-scope="data")
          div.list__tile__title {{ $tc(data.item.text, 2) }}
</template>

<script>
import formBaseMixin from '@/mixins/form/base';

import { PERIODIC_REFRESH_UNITS } from '@/constants';

export default {
  inject: ['$validator'],
  mixins: [formBaseMixin],
  model: {
    prop: 'periodicRefresh',
    event: 'input',
  },
  props: {
    periodicRefresh: {
      type: Object,
      default: () => ({
        enabled: false,
        interval: 0,
        unit: 's',
      }),
    },
    label: {
      type: String,
      required: false,
    },
  },
  computed: {
    availableTypes() {
      return Object.values(PERIODIC_REFRESH_UNITS);
    },
  },
};
</script>
