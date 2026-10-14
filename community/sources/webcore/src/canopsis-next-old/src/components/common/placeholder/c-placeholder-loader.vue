<template functional>
  <div
    :style="{ height: `${props.height}px` }"
    class="placeholder"
  >
    <div class="placeholder__background" />
  </div>
</template>

<script>
export default {
  props: {
    height: {
      type: [String, Number],
      default: 24,
    },
  },
};
</script>

<style lang="scss">
@keyframes placeholder-glow {
  50% {
    opacity: 1;
  }
}

.placeholder {
  --c-placeholder-loader-background: darkgray;

  width: 100%;
  position: absolute;
  top: 0;
  left: 0;
  z-index: 0;

  &__background {
    opacity: .4;
    transform: translate3d(-50%, -50%, 0);
    top: 50%;
    left: 50%;
    height: 86%;
    position: relative;
    background: var(--c-placeholder-loader-background);
  }
}
</style>
