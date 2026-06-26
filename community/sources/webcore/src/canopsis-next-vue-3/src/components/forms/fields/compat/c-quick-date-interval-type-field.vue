<template>
  <c-select
    :model-value="value"
    :items="ranges"
    item-title="text"
    item-value="value"
    v-bind="$attrs"
    @update:model-value="$emit('update:value', $event)"
  />
</template>

<script>
export default {
  emits: ['update:value'],
  props: {
    value: {
      type: String,
      default: '',
    },
    ranges: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
