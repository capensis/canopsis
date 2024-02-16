<template>
  <v-btn-toggle
    class="window-size transparent"
    v-field="value"
    light
    tile
    group
    mandatory
  >
    <v-btn
      v-for="size in availableSizes"
      :ripple="false"
      :key="size.value"
      :value="size.value"
      :dark="dark"
      :light="light"
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
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: WIDGET_GRID_SIZES_KEYS.desktop,
    },
    dark: {
      type: Boolean,
      required: false,
    },
    light: {
      type: Boolean,
      required: false,
    },
    color: {
      type: String,
      required: false,
    },
  },
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
