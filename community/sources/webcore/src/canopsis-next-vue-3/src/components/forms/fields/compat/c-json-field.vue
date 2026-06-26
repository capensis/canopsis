<script setup>
import { computed } from 'vue';

const props = defineProps({
  value: {
    type: String,
    default: '',
  },
});

const emit = defineEmits(['update:value']);

const model = computed({
  get: () => props.value,
  set: value => emit('update:value', value),
});
</script>

<template>
  <c-textarea
    v-bind="$attrs"
    v-model="model"
    rows="8"
  />
</template>
