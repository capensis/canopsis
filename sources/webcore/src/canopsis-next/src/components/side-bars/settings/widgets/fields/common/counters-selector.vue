<template lang="pug">
  v-list-group
    v-list-tile(slot="activator")
      div(:class="validationHeaderClass") {{ $t('settings.counters') }}
    v-container
      v-layout(align-center)
        v-switch(
          v-field="value.enabled",
          color="primary",
          hide-details
        )
        pbehavior-type-field(
          v-field="value.types",
          :disabled="!value.enabled",
          :is-item-disabled="isItemDisabled",
          with-icon,
          chips,
          multiple
        )
</template>

<script>
import { COUNTERS_LIMIT } from '@/constants';

import formValidationHeaderMixin from '@/mixins/form/validation-header';

import PeriodicRefreshField from '@/components/forms/fields/periodic-refresh-field.vue';
import PbehaviorTypeField from '@/components/other/pbehavior/calendar/partials/pbehavior-type-field.vue';

export default {
  inject: ['$validator'],
  components: { PeriodicRefreshField, PbehaviorTypeField },
  mixins: [formValidationHeaderMixin],
  props: {
    value: {
      type: Object,
      default: () => ({ enabled: false, types: [] }),
    },
  },
  methods: {
    isItemDisabled(item) {
      const { types } = this.value;

      return types.length === COUNTERS_LIMIT && !types.includes(item._id);
    },
  },
};
</script>
