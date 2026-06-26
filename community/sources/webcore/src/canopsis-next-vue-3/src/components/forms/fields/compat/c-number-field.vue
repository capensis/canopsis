<script setup>
import { computed } from 'vue';

const props = defineProps({
  value: {
    type: [Number, String],
    default: undefined,
  },
});

const emit = defineEmits(['update:value']);

const model = computed({
  get: () => props.value,
  set: value => emit('update:value', value === '' || value === null ? value : Number(value)),
});
</script>

<template>
  <c-text-field
    v-bind="$attrs"
    v-model="model"
    type="number"
  />
</template>
