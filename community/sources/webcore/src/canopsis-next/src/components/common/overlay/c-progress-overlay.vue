<template>
  <v-fade-transition v-if="transition">
    <v-layout
      v-if="pending"
      class="c-progress-overlay"
      align-center
      justify-center
      column
    >
      <div
        :class="backgroundColor"
        :style="{ opacity }"
        class="c-progress-overlay__background"
      />
      <slot name="progress">
        <v-progress-circular
          :color="color"
          :size="size"
          :width="width"
          indeterminate
        />
      </slot>
      <div v-if="$slots.default" class="c-progress-overlay__content">
        <slot />
      </div>
    </v-layout>
  </v-fade-transition>
  <v-layout
    v-else-if="pending"
    :class="backgroundColor"
    class="c-progress-overlay"
    column
  >
    <div
      :class="backgroundColor"
      :style="{ opacity }"
      class="c-progress-overlay__background"
    />
    <slot name="progress">
      <v-progress-circular
        :color="color"
        :size="size"
        :width="width"
        indeterminate
      />
    </slot>
    <div v-if="$slots.default" class="c-progress-overlay__content">
      <slot />
    </div>
  </v-layout>
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
    size: {
      type: Number,
      required: false,
    },
    width: {
      type: Number,
      required: false,
    },
    backgroundColor: {
      type: String,
      default: 'background',
    },
    color: {
      type: String,
      default: 'primary',
    },
    transition: {
      type: Boolean,
      default: true,
    },
  },
};
</script>

<style lang="scss" scoped>
  .c-progress-overlay {
    position: absolute;
    top: 0;
    left: 0;
    bottom: 0;
    right: 0;
    z-index: 2;

    &__background {
      position: absolute;
      top: 0;
      left: 0;
      bottom: 0;
      right: 0;
    }

    & ::v-deep .v-progress-circular {
      z-index: 3;
    }

    &__content {
      position: relative;
      z-index: 3;
    }
  }
</style>
