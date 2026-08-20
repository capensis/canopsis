<script setup>
import { computed } from 'vue';
import { useTheme } from 'vuetify';

import { isCSSVariable, isDarkColor } from '@/helpers/color';
import { getThemeVariableKeyFromCssVariable } from '@/helpers/entities/theme/entity';

defineOptions({
  inheritAttrs: false,
});

const props = defineProps({
  type: {
    type: String,
    default: undefined,
  },
  color: {
    type: String,
    default: undefined,
  },
});

const theme = useTheme();

const computedColor = computed(() => {
  if (props.color) {
    return props.color;
  }

  if (props.type) {
    return `var(--v-${props.type}-background-base)`;
  }

  return undefined;
});

const iconColor = computed(() => (props.type ? `var(--v-${props.type}-base)` : undefined));

const resolvedBackgroundColor = computed(() => {
  const backgroundColor = computedColor.value;

  if (!backgroundColor) {
    return '#fff';
  }

  if (isCSSVariable(backgroundColor)) {
    const backgroundColorKey = getThemeVariableKeyFromCssVariable(backgroundColor);

    return theme.current.value.colors?.[backgroundColorKey] || '#fff';
  }

  return theme.current.value.colors?.[backgroundColor] || backgroundColor;
});

const isDarkBackground = computed(() => isDarkColor(resolvedBackgroundColor.value));

const alertClasses = computed(() => ({
  'c-alert': true,
  'v-alert--dark-background': isDarkBackground.value,
  'v-alert--light-background': !isDarkBackground.value,
}));

const alertStyles = computed(() => (iconColor.value ? { '--alert-icon-color': iconColor.value } : undefined));
</script>

<template>
  <v-alert
    v-bind="$attrs"
    :type="type"
    :color="computedColor"
    :class="alertClasses"
    :style="alertStyles"
  >
    <slot />
  </v-alert>
</template>

<style lang="scss">
.c-alert.v-alert {
  --alert-color-dark: #212121;
  --alert-color-light: #d9d9d9;

  width: 100%;
  font-weight: 700;
  font-size: 14px;
  line-height: 22px;

  &.v-alert--dark-background {
    color: var(--alert-color-light);
  }

  &.v-alert--light-background {
    color: var(--alert-color-dark);
  }

  .v-alert__prepend .v-icon {
    color: var(--alert-icon-color, currentColor);
  }
}
</style>
