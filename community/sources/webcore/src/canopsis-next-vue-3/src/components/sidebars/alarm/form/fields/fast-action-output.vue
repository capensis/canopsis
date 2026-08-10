<template>
  <widget-settings-item :title="label" optional>
    <c-enabled-field v-model="enabled" />
    <c-text-field
      v-model="output"
      :label="t('common.output')"
      :disabled="!enabled"
    />
  </widget-settings-item>
</template>

<script setup>
import { computed } from 'vue';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';
import { useI18n } from '@/hooks/i18n';

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({}),
  },
  label: {
    type: String,
    default: '',
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

const enabled = computed({
  get: () => props.modelValue.enabled ?? false,
  set: value => updateField('enabled', value),
});

const output = computed({
  get: () => props.modelValue.value ?? '',
  set: value => updateField('value', value),
});
</script>
