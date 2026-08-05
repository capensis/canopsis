<template>
  <c-radio-group
    :model-value="modelValue"
    class="mt-0"
    hide-details
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-radio
      v-for="{ value, text } in availableSizes"
      :key="value"
      :value="value"
      :label="text"
      color="primary"
    />
  </c-radio-group>
</template>

<script>
import { THEME_FONT_SIZES } from '@/constants/theme';

export default {
  emits: ['update:modelValue'],
  props: {
    modelValue: {
      type: Number,
      default: THEME_FONT_SIZES.medium,
    },
  },
  computed: {
    availableSizes() {
      return Object.values(THEME_FONT_SIZES).map(value => ({
        text: this.$t(`theme.fontSize.sizes.${value}`),
        value,
      }));
    },
  },
};
</script>
