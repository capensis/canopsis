<template>
  <v-speed-dial
    v-bind="$attrs"
    :attach="attach"
    :transition="transition"
    class="c-speed-dial"
  >
    <template #activator="{ props: activatorProps, isActive }">
      <slot :props="activatorProps" :is-active="isActive" name="activator" />
    </template>
    <slot />
  </v-speed-dial>
</template>

<script>
export default {
  inheritAttrs: false,

  props: {
    transition: {
      type: String,
      default: 'scale-transition',
    },
    attach: {
      type: [Boolean, String, Object],
      default: true,
    },
  },
};
</script>

<style lang="scss">
.c-speed-dial {
  .v-speed-dial__content {
    padding: 8px 0;
  }
}

.v-speed-dial__content > * > *:nth-child(1) { transition-delay: 0.001s; }
.v-speed-dial__content > * > *:nth-child(2) { transition-delay: 0.05s; }
.v-speed-dial__content > * > *:nth-child(3) { transition-delay: 0.1s; }
.v-speed-dial__content > * > *:nth-child(4) { transition-delay: 0.15s; }
.v-speed-dial__content > * > *:nth-child(5) { transition-delay: 0.2s; }
.v-speed-dial__content > * > *:nth-child(6) { transition-delay: 0.25s; }
.v-speed-dial__content > * > *:nth-child(7) { transition-delay: 0.3s; }
.v-speed-dial__content > * > *:nth-child(8) { transition-delay: 0.35s; }
.v-speed-dial__content > * > *:nth-child(9) { transition-delay: 0.4s; }
.v-speed-dial__content > * > *:nth-child(10) { transition-delay: 0.45s; }

.c-speed-dial__btn {
  .v-btn__content {
    position: relative;
    width: 100%;
    height: 100%;

    > .v-icon:not(:only-child) {
      position: absolute;
      top: 50%;
      left: 50%;
      margin: 0;
      transition: opacity 0.3s cubic-bezier(0.25, 0.8, 0.5, 1),
        transform 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);

      &:first-child {
        opacity: 1;
        transform: translate(-50%, -50%);
      }

      &:last-child {
        opacity: 0;
        transform: translate(-50%, -50%) rotate(-45deg);
      }
    }
  }

  &--active .v-btn__content > .v-icon:not(:only-child) {
    &:first-child {
      opacity: 0;
      transform: translate(-50%, -50%) rotate(45deg);
    }

    &:last-child {
      opacity: 1;
      transform: translate(-50%, -50%);
    }
  }
}
</style>
