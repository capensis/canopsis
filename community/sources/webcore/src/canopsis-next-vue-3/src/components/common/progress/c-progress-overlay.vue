<template>
  <v-overlay
    :model-value="pending"
    :persistent="persistent"
    :transition="transitionValue"
    class="align-center justify-center"
  >
    <v-progress-circular :color="color" :size="size" indeterminate />
  </v-overlay>
</template>

<script>
export default {
  name: 'CProgressOverlay',

  props: {
    pending: {
      type: Boolean,
      default: false,
    },
    persistent: {
      type: Boolean,
      default: true,
    },
    transition: {
      type: [Boolean, String],
      default: 'fade-transition',
    },
    color: {
      type: String,
      default: 'primary',
    },
    size: {
      type: [Number, String],
      default: 64,
    },
  },

  computed: {
    transitionValue() {
      return this.transition || false;
    },
  },
};
</script>
