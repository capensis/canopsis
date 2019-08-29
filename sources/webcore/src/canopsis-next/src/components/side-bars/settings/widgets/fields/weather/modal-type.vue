<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.serviceWeatherModalTypes.title') }}
    v-container
      v-radio-group(:value="value", @change="$emit('input', $event)")
        v-radio(
          :value="$constants.SERVICE_WEATHER_WIDGET_MODAL_TYPES.moreInfo",
          :label="$t('settings.serviceWeatherModalTypes.fields.moreInfo')"
        )
        v-radio(
          :value="$constants.SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList",
          :label="$t('settings.serviceWeatherModalTypes.fields.alarmList')"
        )
        v-radio(
          :value="$constants.SERVICE_WEATHER_WIDGET_MODAL_TYPES.both",
          :label="$t('settings.serviceWeatherModalTypes.fields.both')"
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
