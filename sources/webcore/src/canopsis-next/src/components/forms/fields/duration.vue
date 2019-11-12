<template lang="pug">
  v-layout(row)
    v-flex(xs8)
      v-text-field(
        v-field="value.duration",
        v-validate="'required|numeric|min_value:1'",
        :label="$t('modals.createSnoozeEvent.fields.duration')",
        :error-messages="errors.collect('duration')",
        data-vv-name="duration",
        type="number"
      )
    v-flex(xs4)
      v-select(
        v-field="value.durationType",
        v-validate="'required'",
        :items="availableTypes",
        :error-messages="errors.collect('durationType')",
        data-vv-name="durationType"
      )
        template(slot="selection", slot-scope="data")
          div.input-group__selections__comma {{ $tc(data.item.text, 2) }}
        template(slot="item", slot-scope="data")
          div.list__tile__title {{ $tc(data.item.text, 2) }}
</template>

<script>
import { DURATION_UNITS } from '@/constants';

export default {
  inject: ['$validator'],
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    availableTypes() {
      return Object.values(DURATION_UNITS);
    },
  },
};
</script>
