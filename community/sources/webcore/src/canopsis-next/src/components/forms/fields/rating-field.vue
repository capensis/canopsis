<template>
  <v-rating
    :model-value="modelValue"
    :readonly="readonly"
    :hover="hover"
    class="layout"
    color="success"
    density="comfortable"
    half-increments
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
export default {
  props: {
    modelValue: {
      type: [Number, String],
      default: 0,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    hover: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],
};
</script>
