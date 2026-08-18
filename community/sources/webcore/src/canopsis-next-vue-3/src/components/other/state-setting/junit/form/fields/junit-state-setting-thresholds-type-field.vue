<template>
  <c-radio-group
    :model-value="modelValue"
    :label="label"
    :name="name"
    class="mt-0"
    color="primary"
    hide-details
    mandatory
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-radio :label="t('common.variableTypes.number')" :value="numberValue" color="primary" />
    <v-radio :label="t('common.percent')" :value="percentValue" color="primary" />
  </c-radio-group>
</template>

<script>
import { JUNIT_STATE_SETTING_THRESHOLDS_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: [Number, String],
      default: JUNIT_STATE_SETTING_THRESHOLDS_TYPES.number,
    },
    label: {
      type: String,
      default: '',
    },
    numberValue: {
      type: [Number, String],
      default: JUNIT_STATE_SETTING_THRESHOLDS_TYPES.number,
    },
    percentValue: {
      type: [Number, String],
      default: JUNIT_STATE_SETTING_THRESHOLDS_TYPES.percent,
    },
    name: {
      type: String,
      default: 'type',
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    return { t };
  },
};
</script>
