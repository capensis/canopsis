<template lang="pug">
  v-btn-toggle.window-size.transparent(v-field="value", mandatory, light, dense)
    v-btn(
      v-for="size in availableSizes",
      :ripple="false",
      :key="size.value",
      :value="size.value",
      icon
    )
      v-icon {{ size.icon }}
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
