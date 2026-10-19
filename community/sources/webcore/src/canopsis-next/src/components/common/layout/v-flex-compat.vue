<template>
  <div v-bind="attrs" :class="widthClass" class="v-flex-compat">
    <slot />
  </div>
</template>

<script setup>
import { computed, useAttrs } from 'vue';

defineOptions({
  inheritAttrs: false,
});

const props = defineProps({
  xs1: {
    type: Boolean,
    default: false,
  },
  xs2: {
    type: Boolean,
    default: false,
  },
  xs3: {
    type: Boolean,
    default: false,
  },
  xs4: {
    type: Boolean,
    default: false,
  },
  xs5: {
    type: Boolean,
    default: false,
  },
  xs6: {
    type: Boolean,
    default: false,
  },
  xs7: {
    type: Boolean,
    default: false,
  },
  xs8: {
    type: Boolean,
    default: false,
  },
  xs9: {
    type: Boolean,
    default: false,
  },
  xs10: {
    type: Boolean,
    default: false,
  },
  xs11: {
    type: Boolean,
    default: false,
  },
  xs12: {
    type: Boolean,
    default: false,
  },
});

const attrs = useAttrs();

const widthClass = computed(() => {
  const size = {
    [props.xs1]: 1,
    [props.xs2]: 2,
    [props.xs3]: 3,
    [props.xs4]: 4,
    [props.xs5]: 5,
    [props.xs6]: 6,
    [props.xs7]: 7,
    [props.xs8]: 8,
    [props.xs9]: 9,
    [props.xs10]: 10,
    [props.xs11]: 11,
    [props.xs12]: 12,
  }.true;

  return size ? `v-flex-compat--xs${size}` : '';
});
</script>

<style lang="scss" scoped>
.v-flex-compat {
  flex: 1 1 auto;
  min-width: 0;

  @for $size from 1 through 12 {
    &--xs#{$size} {
      flex: 0 0 calc(100% * #{$size} / 12);
      max-width: calc(100% * #{$size} / 12);
    }
  }
}
</style>
