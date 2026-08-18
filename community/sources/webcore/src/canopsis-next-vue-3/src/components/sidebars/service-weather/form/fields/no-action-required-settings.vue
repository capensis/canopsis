<template>
  <widget-settings-item :title="t('settings.noActionRequiredSettings.title')">
    <c-layout column>
      <c-enabled-color-picker-field
        :model-value="modelValue.color"
        :label="t('settings.noActionRequiredSettings.customColor')"
        column
        @update:model-value="updateField('color', $event)"
      />
      <c-enabled-icon-field
        :model-value="modelValue.icon_name"
        :label="t('settings.noActionRequiredSettings.customIcon')"
        @update:model-value="updateField('icon_name', $event)"
      />
    </c-layout>
  </widget-settings-item>
</template>

<script setup>
import { useI18n } from '@/hooks/i18n';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
});

const emit = defineEmits(['update:modelValue']);

const { t } = useI18n();

const updateField = (field, value) => {
  emit('update:modelValue', {
    ...props.modelValue,
    [field]: value,
  });
};
</script>
