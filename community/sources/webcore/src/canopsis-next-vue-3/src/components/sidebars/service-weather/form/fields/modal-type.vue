<template>
  <widget-settings-item :title="t('settings.serviceWeatherModalTypes.title')">
    <v-radio-group
      :model-value="modelValue"
      @update:model-value="$emit('update:modelValue', $event)"
    >
      <v-radio
        :value="SERVICE_WEATHER_WIDGET_MODAL_TYPES.moreInfo"
        :label="t('settings.serviceWeatherModalTypes.fields.moreInfo')"
        color="primary"
      />
      <v-radio
        :value="SERVICE_WEATHER_WIDGET_MODAL_TYPES.alarmList"
        :label="t('settings.serviceWeatherModalTypes.fields.alarmList')"
        color="primary"
      />
      <v-radio
        :value="SERVICE_WEATHER_WIDGET_MODAL_TYPES.both"
        :label="t('settings.serviceWeatherModalTypes.fields.both')"
        color="primary"
      />
    </v-radio-group>
  </widget-settings-item>
</template>

<script setup>
import { SERVICE_WEATHER_WIDGET_MODAL_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

defineProps({
  modelValue: {
    type: String,
    default: SERVICE_WEATHER_WIDGET_MODAL_TYPES.moreInfo,
  },
});

defineEmits(['update:modelValue']);

const { t } = useI18n();
</script>
