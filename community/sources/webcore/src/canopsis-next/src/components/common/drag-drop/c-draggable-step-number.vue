<template>
  <v-layout
    align-center
    justify-space-around
  >
    <span class="handler">
      <v-icon
        class="draggable"
        v-show="!disabled"
        :class="dragClass"
      >
        drag_indicator
      </v-icon>
    </span>
    <v-avatar
      class="white--text"
      :color="color"
      size="32"
    >
      <slot />
    </v-avatar>
  </v-layout>
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
  },
};
</script>

<style lang="scss">
  .handler {
    height: 24px !important;
    width: 18px !important;
  }
</style>
