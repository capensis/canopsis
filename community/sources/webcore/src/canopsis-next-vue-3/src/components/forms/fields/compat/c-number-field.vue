<template>
  <c-text-field v-bind="$attrs" v-model="model" type="number" />
</template>

<script setup>
import { computed, useAttrs } from 'vue';

const props = defineProps({
  value: {
    type: [Number, String],
    default: undefined,
  },
});

const emit = defineEmits(['update:value']);
const attrs = useAttrs();

const minValue = computed(() => {
  const min = attrs.min;

  return min === null || min === undefined || min === '' ? undefined : Number(min);
});

const normalizeValue = (value) => {
  if (value === '' || value === null) {
    return value;
  }

  const number = Number(value);

  if (Number.isFinite(minValue.value) && minValue.value >= 0 && number < minValue.value) {
    return minValue.value;
  }

  return number;
};

const model = computed({
  get: () => props.value,
  set: value => emit('update:value', normalizeValue(value)),
});
</script>
