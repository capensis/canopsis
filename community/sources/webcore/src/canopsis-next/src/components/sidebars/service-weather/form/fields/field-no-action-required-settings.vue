<template>
  <widget-settings-item :title="$t('settings.noActionRequiredSettings.title')">
    <v-layout column>
      <c-enabled-color-picker-field
        v-field="value.color"
        :label="$t('settings.noActionRequiredSettings.customColor')"
        column
      />
      <c-enabled-icon-field
        v-field="value.icon_name"
        :label="$t('settings.noActionRequiredSettings.customIcon')"
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
