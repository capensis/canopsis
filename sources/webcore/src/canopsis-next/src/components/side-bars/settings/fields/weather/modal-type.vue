<template lang="pug">
  v-list-group(data-test="modalType")
    v-list-tile(slot="activator") {{ $t('settings.serviceWeatherModalTypes.title') }}
    v-container
      v-radio-group(data-test="modalTypeGroup", :value="value", @change="$emit('input', $event)")
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

export default {
  props: {
    value: {
      type: String,
      default: SERVICE_WEATHER_WIDGET_MODAL_TYPES.moreInfo,
    },
  },
};
</script>
