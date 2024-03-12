<template>
  <widget-settings-item :title="$t('settings.colorIndicator.title')">
    <c-color-indicator-field v-field="value" />
  </widget-settings-item>
</template>

<script>
import { COLOR_INDICATOR_TYPES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: String,
      default: COLOR_INDICATOR_TYPES.state,
    },
  },
};
</script>
