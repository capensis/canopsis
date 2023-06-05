<template lang="pug">
  widget-settings-item(:title="$t('settings.serviceWeatherModalTypes.title')")
    v-radio-group(v-field="value", :name="name")
      v-radio(
        :value="$constants.SERVICE_WEATHER_WIDGET_MODAL_TYPES.moreInfo",
        :label="$t('settings.serviceWeatherModalTypes.fields.moreInfo')",
        color="primary"
      )
      v-radio(
        :value="$constants.SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList",
        :label="$t('settings.serviceWeatherModalTypes.fields.alarmList')",
        color="primary"
      )
      v-radio(
        :value="$constants.SERVICE_WEATHER_WIDGET_MODAL_TYPES.both",
        :label="$t('settings.serviceWeatherModalTypes.fields.both')",
        color="primary"
      )
</template>

<script>
import { SERVICE_WEATHER_WIDGET_MODAL_TYPES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: String,
      default: SERVICE_WEATHER_WIDGET_MODAL_TYPES.moreInfo,
    },
    name: {
      type: String,
      default: 'modalType',
    },
  },
};
</script>
