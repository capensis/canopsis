<template>
  <v-icon v-bind="iconAttrs" :class="iconClass">
    <slot />
  </v-icon>
</template>

<script setup>
import { computed, useAttrs } from 'vue';

defineOptions({
  inheritAttrs: false,
});

const props = defineProps({
  noFill: {
    type: Boolean,
    default: false,
  },
  size: {
    type: [Number, String],
    default: undefined,
  },
  xSmall: {
    type: Boolean,
    default: false,
  },
  small: {
    type: Boolean,
    default: false,
  },
  large: {
    type: Boolean,
    default: false,
  },
  xLarge: {
    type: Boolean,
    default: false,
  },
});

const attrs = useAttrs();

const legacySize = computed(
  () => ({
      [props.xSmall]: 'x-small',
      [props.small]: 'small',
      [props.large]: 'large',
      [props.xLarge]: 'x-large',
    }).true,
);

const iconAttrs = computed(() => ({
  ...attrs,
  size: props.size ?? legacySize.value,
}));

const iconClass = computed(() => ({
  'v-icon--no-fill': props.noFill,
}));
</script>
