<template>
  <widget-settings-item :title="t('settings.counters')">
    <c-layout column>
      <c-enabled-field
        :model-value="modelValue.pbehavior_enabled"
        :label="t('settings.pbehaviorCounters')"
        hide-details
        @update:model-value="updateField('pbehavior_enabled', $event)"
      />
      <c-pbehavior-type-field
        :model-value="modelValue.pbehavior_types"
        :required="modelValue.pbehavior_enabled"
        :disabled="!modelValue.pbehavior_enabled"
        chips
        multiple
        @update:model-value="updateField('pbehavior_types', $event)"
      />
      <c-enabled-field
        :model-value="modelValue.state_enabled"
        :label="t('settings.entityStateCounters')"
        hide-details
        @update:model-value="updateField('state_enabled', $event)"
      />
      <c-service-weather-state-counter-field
        :model-value="modelValue.state_types"
        :required="modelValue.state_enabled"
        :disabled="!modelValue.state_enabled"
        @update:model-value="updateField('state_types', $event)"
      />
    </c-layout>
  </widget-settings-item>
</template>

<script setup>
import { onMounted } from 'vue';

import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
});

const emit = defineEmits(['update:modelValue']);

const { t } = useI18n();
const { fetchPbehaviorTypesFieldList } = usePbehaviorType();

const updateField = (field, value) => {
  emit('update:modelValue', {
    ...props.modelValue,
    [field]: value,
  });
};

onMounted(() => {
  fetchPbehaviorTypesFieldList({
    params: {
      paginate: false,
      with_hidden: true,
      types: [
        PBEHAVIOR_TYPE_TYPES.inactive,
        PBEHAVIOR_TYPE_TYPES.maintenance,
        PBEHAVIOR_TYPE_TYPES.pause,
      ],
    },
  });
});
</script>
