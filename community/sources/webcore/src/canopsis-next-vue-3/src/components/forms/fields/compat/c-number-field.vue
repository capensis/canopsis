<template>
  <c-text-field
    v-bind="$attrs"
    v-model="model"
    :min="min"
    :max="max"
    type="number"
  />
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  value: {
    type: [Number, String],
    default: undefined,
  },
  min: {
    type: Number,
    default: undefined,
  },
  max: {
    type: Number,
    default: undefined,
  },
});

/**
 * Converts a bound min/max constraint to a number, leaving empty values as-is.
 *
 * @param {number|string|null|undefined} value Constraint from props.
 * @returns {number|string|null|undefined}
 */
const normalizeNumber = (value) => {
  if (value === '' || value === null || value === undefined) {
    return value;
  }

  return Number(value);
};

const emit = defineEmits(['update:value']);

/**
 * Numeric minimum bound used for clamping. Empty when `min` is not set.
 *
 * @type {import('vue').ComputedRef<number|string|null|undefined>}
 */
const minValue = computed(() => normalizeNumber(props.min));

/**
 * Numeric maximum bound used for clamping. Empty when `max` is not set.
 *
 * @type {import('vue').ComputedRef<number|string|null|undefined>}
 */
const maxValue = computed(() => normalizeNumber(props.max));

/**
 * Coerces the input to a number and clamps it to `min`/`max` when those bounds are finite.
 *
 * Empty string and `null` are kept so the field can be cleared.
 *
 * @param {number|string|null} value Raw input value.
 * @returns {number|string|null}
 */
const normalizeValue = (value) => {
  if (value === '' || value === null) {
    return value;
  }

  const number = Number(value);

  if (Number.isFinite(minValue.value) && number < minValue.value) {
    return minValue.value;
  }

  if (Number.isFinite(maxValue.value) && number > maxValue.value) {
    return maxValue.value;
  }

  return number;
};

/**
 * Two-way binding between the parent `value` and the inner text field.
 *
 * @type {import('vue').WritableComputedRef<number|string|null|undefined>}
 */
const model = computed({
  get: () => props.value,
  set: value => emit('update:value', normalizeValue(value)),
});
</script>
