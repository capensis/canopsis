<template lang="pug">
  widget-settings-item(:title="$t('settings.colorIndicator.title')")
    c-color-indicator-field(v-field="value")
</template>

<script>
import { COLOR_INDICATOR_TYPES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

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
