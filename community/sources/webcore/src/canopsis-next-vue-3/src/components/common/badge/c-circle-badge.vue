<template>
  <span
    :class="{ 'c-circle-badge--outlined': outlined }"
    :style="badgeStyle"
    class="c-circle-badge"
  >
    <slot />
  </span>
</template>

<script>
export default {
  name: 'CCircleBadge',

  props: {
    color: {
      type: String,
      default: 'primary',
    },
    outlined: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    badgeStyle() {
      const color = this.color === 'white'
        ? 'white'
        : `rgb(var(--v-theme-${this.color}))`

      return {
        '--c-circle-badge-color': color,
      }
    },
  },
}
</script>

<style lang="scss" scoped>
.c-circle-badge {
  align-items: center;
  background: var(--c-circle-badge-color);
  border: 1px solid var(--c-circle-badge-color);
  border-radius: 50%;
  color: rgb(var(--v-theme-on-primary));
  display: inline-flex;
  font-size: 12px;
  height: 22px;
  justify-content: center;
  min-width: 22px;
  padding: 0 4px;

  &--outlined {
    background: transparent;
  }
}
</style>
