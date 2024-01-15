<template>
  <widget-settings-item :title="$t('settings.actionRequiredSettings.title')">
    <v-layout column>
      <c-enabled-field
        v-field="value.is_blinking"
        :label="$t('settings.actionRequiredSettings.isBlinking')"
      />
      <c-enabled-color-picker-field
        v-field="value.color"
        :label="$t('settings.actionRequiredSettings.customColor')"
        column
      />
      <c-enabled-icon-field
        v-field="value.icon_name"
        :label="$t('settings.actionRequiredSettings.customIcon')"
      />
    </v-layout>
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Object,
      required: false,
    },
  },
};
</script>
