<template>
  <v-radio-group
    class="mt-0"
    v-field="value"
    hide-details
  >
    <v-radio
      v-for="{ value: size, text } in availableSizes"
      :key="size"
      :value="size"
      :label="text"
      color="primary"
    />
  </v-radio-group>
</template>

<script>
import { THEME_FONT_SIZES } from '@/constants/theme';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
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
