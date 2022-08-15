<template lang="pug">
  v-icon.mermaid-point(
    v-on="$listeners",
    :class="{ 'mermaid-point--linked': !entity }",
    :style="{ top: `${y}px`, left: `${x}px` }",
    :size="size",
    color="grey darken-2"
    ) {{ icon }}
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
    entity: {
      type: [String, Object],
      required: false,
    },
    size: {
      type: Number,
      default: 24,
    },
  },
  computed: {
    icon() {
      return this.entity ? 'location_on' : 'link';
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

  &--linked {
    transform: translate(-50%, -50%);
  }
}
</style>
