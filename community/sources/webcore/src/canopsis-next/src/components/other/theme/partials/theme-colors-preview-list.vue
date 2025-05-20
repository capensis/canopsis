<template>
  <v-layout class="gap-4" column>
    <span class="text-subtitle-2 font-italic">{{ $t('theme.checkColors') }}:</span>
    <theme-colors-preview
      v-for="(item, index) in items"
      :key="index"
      :color="item.color"
      :background-color="item.backgroundColor"
      :font-size="item.fontSize"
      :text="item.text"
    />
  </v-layout>
</template>

<script>
import ThemeColorsPreview from './theme-colors-preview.vue';

export default {
  components: { ThemeColorsPreview },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
