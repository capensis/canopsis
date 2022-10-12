<template lang="pug">
  widget-settings-item(:title="$t('settings.counters')")
    v-layout(align-center)
      v-switch(
        v-field="value.pbehavior_enabled",
        color="primary",
        hide-details
      )
      c-pbehavior-type-field(
        v-field="value.pbehavior_types",
        :required="value.pbehavior_enabled",
        :disabled="!value.pbehavior_enabled",
        :is-item-disabled="isItemDisabled",
        with-icon,
        chips,
        multiple
      )
    v-layout(align-center)
      v-switch(
        v-field="value.state_enabled",
        color="primary",
        hide-details
      )
      c-service-weather-state-counter-field(
        v-field="value.state_types",
        :required="value.state_enabled",
        :disabled="!value.state_enabled"
      )
</template>

<script>
import { COUNTERS_LIMIT } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Object,
      required: false,
    },
  },
  methods: {
    isItemDisabled(item) {
      const { pbehavior_types: types } = this.value;

      return types.length === COUNTERS_LIMIT && !types.includes(item._id);
    },
  },
};
</script>
