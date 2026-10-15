<template>
  <div :style="{ height: `${height}px` }" class="c-placeholder-loader">
    <div class="c-placeholder-loader__background" />
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
@keyframes c-placeholder-loader-glow {
  50% {
    opacity: 1;
  }
}

.c-placeholder-loader {
  --c-placeholder-loader-background: darkgray;

  position: absolute;
  top: 0;
  left: 0;
  z-index: 0;
  width: 100%;

  &__background {
    position: relative;
    top: 50%;
    left: 50%;
    height: 86%;
    background: var(--c-placeholder-loader-background);
    opacity: .4;
    transform: translate3d(-50%, -50%, 0);
  }
}
</style>
