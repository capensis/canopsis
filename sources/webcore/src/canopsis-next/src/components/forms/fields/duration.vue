<template lang="pug">
  v-layout(row)
    v-flex(xs8)
      v-text-field(
        type="number",
        :label="$t('modals.createSnoozeEvent.fields.duration')",
        :error-messages="errors.collect('duration')",
        :value="value.duration",
        @input="updateField('duration', $event)",
        v-validate="'required|numeric|min_value:1'",
        data-vv-name="duration"
      )
    v-flex(xs4)
      v-select(
        :items="availableTypes",
        :value="value.durationType",
        @input="updateField('durationType', $event)",
        v-validate="'required'",
        data-vv-name="durationType",
        :error-messages="errors.collect('durationType')"
      )
        template(slot="selection", slot-scope="data")
          div.input-group__selections__comma {{ $tc(data.item.text, 2) }}
        template(slot="item", slot-scope="data")
          div.list__tile__title {{ $tc(data.item.text, 2) }}
</template>

<script>
import { DURATION_UNITS } from '@/constants';

import formMixin from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
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
