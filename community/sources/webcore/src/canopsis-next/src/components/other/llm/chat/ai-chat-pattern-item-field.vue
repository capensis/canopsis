<template>
  <div>
    <c-select
      :model-value="modelValue"
      :label="label"
      :items="items"
      name="patternItem"
      item-title="text"
      item-value="value"
      hide-details
      required
      return-object
      @update:model-value="$emit('update:modelValue', $event)"
    />
  </div>
</template>

<script>
export default {
  emits: ['update:modelValue'],
  props: {
    modelValue: {
      type: Object,
      default: null,
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
  },
};
</script>
