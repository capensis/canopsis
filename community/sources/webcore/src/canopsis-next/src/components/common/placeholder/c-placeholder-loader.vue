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
  width: 100%;
  position: absolute;
  top: 0;
  left: 0;
  z-index: 1;

  &__background {
    opacity: .2;
    transform: translate3d(-50%, -50%, 0);
    top: 50%;
    left: 50%;
    height: 86%;
    position: relative;
    animation: placeholder-glow 2s ease-in-out infinite;
    animation-iteration-count: 2;

    .v-application.theme--light & {
      background: darkgray; // TODO: use vars
    }

    .v-application.theme--dark & {
      background: darkgray; // TODO: use vars
    }
  }
}
</style>
