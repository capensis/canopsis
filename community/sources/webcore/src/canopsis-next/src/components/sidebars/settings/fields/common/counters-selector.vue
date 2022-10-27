<template lang="pug">
  widget-settings-item(:title="$t('settings.counters')")
    v-layout(column)
      c-enabled-field(
        v-field="value.pbehavior_enabled",
        :label="$t('settings.pbehaviorCounters')",
        hide-details
      )
      c-pbehavior-type-field(
        v-field="value.pbehavior_types",
        :required="value.pbehavior_enabled",
        :disabled="!value.pbehavior_enabled",
        :max="$constants.PBEHAVIOR_COUNTERS_LIMIT",
        :types="pbehaviorTypeTypes",
        chips,
        multiple
      )
      c-enabled-field(
        v-field="value.state_enabled",
        :label="$t('settings.entityStateCounters')",
        hide-details
      )
      c-service-weather-state-counter-field(
        v-field="value.state_types",
        :required="value.state_enabled",
        :disabled="!value.state_enabled"
      )
</template>

<script>
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Object,
      required: false,
    },
  },
  computed: {
    pbehaviorTypeTypes() {
      return [PBEHAVIOR_TYPE_TYPES.inactive, PBEHAVIOR_TYPE_TYPES.maintenance, PBEHAVIOR_TYPE_TYPES.pause];
    },
  },
};
</script>
