<template>
  <v-fade-transition v-if="transition">
    <div v-if="pending" class="c-progress-overlay">
      <div :style="backgroundStyle" class="c-progress-overlay__background" />

      <slot name="progress">
        <v-progress-circular :color="color" :size="size" :width="width" indeterminate />
      </slot>

      <div v-if="$slots.default" class="c-progress-overlay__content">
        <slot />
      </div>
    </div>
  </v-fade-transition>

  <div v-else-if="pending" class="c-progress-overlay">
    <div :style="backgroundStyle" class="c-progress-overlay__background" />

    <slot name="progress">
      <v-progress-circular :color="color" :size="size" :width="width" indeterminate />
    </slot>

    <div v-if="$slots.default" class="c-progress-overlay__content">
      <slot />
    </div>
  </div>
</template>

<script>
export default {
  name: 'CProgressOverlay',

  props: {
    pending: {
      type: Boolean,
      default: false,
    },
    opacity: {
      type: Number,
      default: 0.5,
    },
    transition: {
      type: Boolean,
      default: true,
    },
    backgroundColor: {
      type: String,
      default: 'background',
    },
    color: {
      type: String,
      default: 'primary',
    },
    size: {
      type: [Number, String],
      required: false,
    },
    width: {
      type: [Number, String],
      required: false,
    },
  },

  computed: {
    backgroundStyle() {
      return {
        opacity: this.opacity,
        backgroundColor: `rgb(var(--v-theme-${this.backgroundColor}))`,
      };
    },
  },
};
</script>

<style lang="scss" scoped>
.c-progress-overlay {
  position: absolute;
  inset: 0;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  &__background {
    position: absolute;
    inset: 0;
  }

  :deep(.v-progress-circular) {
    z-index: 3;
  }

  &__content {
    position: relative;
    z-index: 3;
  }
}
</style>
