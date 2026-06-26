<script setup>
import { computed, useAttrs } from 'vue';

defineOptions({
  inheritAttrs: false,
});

const props = defineProps({
  tag: {
    type: String,
    default: 'div',
  },
  column: {
    type: Boolean,
    default: false,
  },
  row: {
    type: Boolean,
    default: false,
  },
  inline: {
    type: Boolean,
    default: false,
  },
  wrap: {
    type: Boolean,
    default: false,
  },
  fill: {
    type: Boolean,
    default: false,
  },
  align: {
    type: String,
    default: '',
  },
  justify: {
    type: String,
    default: '',
  },
  gap: {
    type: [Number, String],
    default: null,
  },
  rowGap: {
    type: [Number, String],
    default: null,
  },
  columnGap: {
    type: [Number, String],
    default: null,
  },
});

const attrs = useAttrs();

const classList = value => {
  if (!value) {
    return [];
  }

  if (Array.isArray(value)) {
    return value.flatMap(classList);
  }

  if (typeof value === 'object') {
    return Object.entries(value)
      .filter(([, enabled]) => enabled)
      .map(([className]) => className);
  }

  return String(value).split(/\s+/).filter(Boolean);
};

const normalizeGapClass = className => className.replace(/^gap-/, 'ga-');
const toGapClass = (prefix, value) => `${prefix}-${value}`;

const normalizedAttrClasses = computed(() => (
  classList(attrs.class).map(normalizeGapClass)
));

const layoutClasses = computed(() => [
  normalizedAttrClasses.value,
  props.inline ? 'd-inline-flex' : 'd-flex',
  props.column ? 'flex-column align-stretch' : 'flex-row',
  {
    'flex-wrap': props.wrap,
    'c-layout--fill': props.fill,
  },
  props.align && `align-${props.align}`,
  props.justify && `justify-${props.justify}`,
  props.gap !== null && toGapClass('ga', props.gap),
  props.rowGap !== null && toGapClass('gr', props.rowGap),
  props.columnGap !== null && toGapClass('gc', props.columnGap),
]);

const layoutAttrs = computed(() => {
  const rest = { ...attrs };

  delete rest.class;

  return rest;
});
</script>

<template>
  <component
    :is="tag"
    v-bind="layoutAttrs"
    :class="layoutClasses"
  >
    <slot />
  </component>
</template>

<style lang="scss" scoped>
.c-layout--fill {
  flex: 1 1 auto;
  min-height: 0;
  min-width: 0;
}
</style>
