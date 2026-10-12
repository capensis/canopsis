<template>
  <widget-settings-item :title="t('settings.actionRequiredSettings.title')">
    <c-layout column>
      <c-enabled-field
        :model-value="modelValue.is_blinking"
        :label="t('settings.actionRequiredSettings.isBlinking')"
        @update:model-value="updateField('is_blinking', $event)"
      />
      <c-enabled-color-picker-field
        :model-value="modelValue.color"
        :label="t('settings.actionRequiredSettings.customColor')"
        name="actionRequiredSettings.color"
        column
        @update:model-value="updateField('color', $event)"
      />
      <c-enabled-icon-field
        :model-value="modelValue.icon_name"
        :label="t('settings.actionRequiredSettings.customIcon')"
        name="actionRequiredSettings.icon_name"
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
