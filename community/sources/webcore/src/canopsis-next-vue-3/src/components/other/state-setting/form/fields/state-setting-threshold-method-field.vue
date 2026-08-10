<template>
  <c-radio-group
    :model-value="modelValue"
    :label="label"
    :disabled="disabled"
    hide-details
    inline
    mandatory
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-radio
      v-for="method in methods"
      :key="method"
      :label="t(`stateSetting.thresholdMethods.${method}`)"
      :value="method"
      color="primary"
    />
  </c-radio-group>
</template>

<script>
import { STATE_SETTING_THRESHOLDS_METHODS } from '@/constants';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      default: STATE_SETTING_THRESHOLDS_METHODS.share,
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { t } = useI18n();

    return {
      t,
      methods: STATE_SETTING_THRESHOLDS_METHODS,
    };
  },
};
</script>
