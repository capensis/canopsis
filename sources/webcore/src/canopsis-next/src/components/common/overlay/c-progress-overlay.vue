<template lang="pug">
  v-fade-transition
    v-layout.progress(
      v-show="pending",
      :class="backgroundColor",
      :style="{ opacity }",
      data-test="progressOverlay"
    )
      v-progress-circular(indeterminate, color="primary")
</template>

<script>
export default {
  props: {
    pending: {
      type: Boolean,
      default: false,
    },
    opacity: {
      type: Number,
      default: 0.5,
    },
    backgroundColor: {
      type: String,
      default: 'white',
    },
  },
};
</script>

<style lang="scss" scoped>
  .progress {
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    right: 0;
    z-index: 2;

    & /deep/ .v-progress-circular {
      top: 50%;
      left: 50%;
      margin-top: -16px;
      margin-left: -16px;
    }
  }
</style>
