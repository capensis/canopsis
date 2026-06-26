<template>
  <v-btn
    :loading="loading"
    size="small"
    variant="text"
    icon
    @click="$emit('update:modelValue', !modelValue)"
  >
    <c-icon :color="color">
      {{ modelValue ? 'keyboard_arrow_up' : 'keyboard_arrow_down' }}
    </c-icon>
  </v-btn>
</template>

<script>
export default {
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    color: {
      type: String,
      required: false,
    },
    loading: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],
};
</script>
