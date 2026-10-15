<template>
  <v-btn-toggle
    :model-value="modelValue"
    class="window-size transparent"
    variant="text"
    mandatory
    group
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <v-btn
      v-for="size in availableSizes"
      :key="size.value"
      :ripple="false"
      :value="size.value"
      icon
    >
      <v-icon :color="color">
        {{ size.icon }}
      </v-icon>
    </v-btn>
  </v-btn-toggle>
</template>

<script>
import { WIDGET_GRID_SIZES_STYLES, WIDGET_GRID_SIZES_KEYS } from '@/constants';

export default {
  props: {
    modelValue: {
      type: String,
      default: WIDGET_GRID_SIZES_KEYS.desktop,
    },
    color: {
      type: String,
      required: false,
    },
  },
  emits: ['update:modelValue'],
  computed: {
    availableSizes() {
      return [
        WIDGET_GRID_SIZES_STYLES.mobile,
        WIDGET_GRID_SIZES_STYLES.tablet,
        WIDGET_GRID_SIZES_STYLES.desktop,
      ];
    },
  },
};
</script>

<style lang="scss" scoped>
.window-size {
  border: none;
}
</style>
