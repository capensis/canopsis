<template>
  <v-btn
    :loading="loading"
    :size="size"
    variant="text"
    icon
    @click="$emit('update:modelValue', !modelValue)"
  >
    <c-icon :color="color" :size="iconSize">
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
    size: {
      type: [Number, String],
      default: 'small',
    },
    iconSize: {
      type: [Number, String],
      default: undefined,
    },
  },

  emits: ['update:modelValue'],
};
</script>
