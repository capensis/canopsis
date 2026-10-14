<template>
  <c-layout align="center" justify="space-around">
    <span class="handler">
      <v-icon v-show="!disabled" :class="dragClass" class="draggable"> drag_indicator </v-icon>
    </span>
    <v-avatar
      v-if="!hideNumber"
      :color="color"
      class="white--text"
      size="32"
    >
      <slot />
    </v-avatar>
  </c-layout>
</template>

<script>
export default {
  props: {
    dragClass: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: 'primary',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    hideNumber: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss">
.handler {
  --draggable-step-number-handler-size: 24px;

  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: var(--draggable-step-number-handler-size) !important;
  height: var(--draggable-step-number-handler-size) !important;

  .v-icon {
    width: var(--draggable-step-number-handler-size);
    height: var(--draggable-step-number-handler-size);
    font-size: var(--draggable-step-number-handler-size);
  }
}
</style>
