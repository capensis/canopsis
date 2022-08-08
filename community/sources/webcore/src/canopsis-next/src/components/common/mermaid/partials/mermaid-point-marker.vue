<template lang="pug">
  v-icon.mermaid-point(
    v-on="$listeners",
    :style="{ top: `${y}px`, left: `${x}px` }",
    :size="size",
    color="grey darken-2"
    ) location_on
</template>

<script>
export default {
  props: {
    x: {
      type: Number,
      required: true,
    },
    y: {
      type: Number,
      required: true,
    },
    size: {
      type: Number,
      default: 24,
    },
  },
};
</script>

<style lang="scss">
.mermaid-point {
  position: absolute;
  transform: translate(-50%, -100%);
  user-select: none;
  cursor: pointer;
  transition: none;
}
</style>
