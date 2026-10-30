<template>
  <v-rating
    :model-value="3.5"
    :readonly="readonly"
    :hover="hover"
    class="layout"
    color="success"
    density="comfortable"
    half-increments
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <template #item="{ props: itemProps }">
      <v-btn v-bind="itemProps">
        <c-icon :icon="itemProps.icon" :no-fill="itemProps.icon === '$ratingEmpty'" />
      </v-btn>
    </template>
  </v-rating>
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
