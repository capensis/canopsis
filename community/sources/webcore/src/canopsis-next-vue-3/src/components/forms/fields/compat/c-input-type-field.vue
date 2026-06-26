<template>
  <c-select
    :model-value="value"
    :items="types"
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
    types: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
