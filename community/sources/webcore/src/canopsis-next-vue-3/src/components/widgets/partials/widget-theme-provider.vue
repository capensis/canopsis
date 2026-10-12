<template>
  <slot />
</template>

<script setup>
import { computed, inject, provide, readonly } from 'vue';

import { WIDGET_COLORS_MODES } from '@/constants';

import { themeTableColorsToForm, themeWidgetColorsToForm } from '@/helpers/entities/theme/form';

const props = defineProps({
  widget: {
    type: Object,
    required: true,
  },
});

const parentSystem = inject('$system', {});

const theme = computed(() => {
  const parentTheme = parentSystem.theme ?? {};
  const widgetColors = props.widget.parameters?.colors;

  if (widgetColors?.mode !== WIDGET_COLORS_MODES.custom) {
    return parentTheme;
  }

  return {
    ...parentTheme,
    colors: {
      ...parentTheme.colors,
      table: themeTableColorsToForm(widgetColors.table),
      widget: themeWidgetColorsToForm(widgetColors.widget),
    },
  };
});

const localSystem = readonly({
  get dark() {
    return parentSystem.dark;
  },
  get timezone() {
    return parentSystem.timezone;
  },
  get dateFormat() {
    return parentSystem.dateFormat;
  },
  get theme() {
    return theme.value;
  },
  setTheme: (...args) => parentSystem.setTheme?.(...args),
});

provide('$system', localSystem);
</script>
