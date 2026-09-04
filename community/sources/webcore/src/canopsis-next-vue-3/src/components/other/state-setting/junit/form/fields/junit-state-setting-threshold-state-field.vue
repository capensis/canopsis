<template>
  <c-layout class="junit-state-setting-threshold-state-field" align="center" gap="2">
    <span class="junit-state-setting-threshold-state-field__label mt-4">
      <span class="text-error">*</span>
      {{ label }}
    </span>
    <c-number-field
      :value="modelValue"
      :name="name"
      :min="0"
      :max="max"
      :suffix="suffix"
      class="junit-state-setting-threshold-state-field__input"
      hide-details="auto"
      required
      @update:value="$emit('update:modelValue', $event)"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { JUNIT_STATE_SETTING_THRESHOLDS_TYPES } from '@/constants';

export default {
  props: {
    modelValue: {
      type: [Number, String],
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    type: {
      type: Number,
      default: JUNIT_STATE_SETTING_THRESHOLDS_TYPES.number,
    },
    name: {
      type: String,
      default: 'state',
    },
  },

  emits: ['update:modelValue'],

  setup(props) {
    const isPercentType = computed(
      () => props.type === JUNIT_STATE_SETTING_THRESHOLDS_TYPES.percent,
    );
    const max = computed(() => (isPercentType.value ? 100 : undefined));
    const suffix = computed(() => (isPercentType.value ? '%' : undefined));

    return {
      max,
      suffix,
    };
  },
};
</script>

<style lang="scss" scoped>
.junit-state-setting-threshold-state-field {
  &__label {
    display: flex;
    align-items: center;
    min-width: 88px;
    white-space: nowrap;
    line-height: 1;
  }

  &__input {
    width: 140px;
    flex: 0 0 auto;
  }
}
</style>
