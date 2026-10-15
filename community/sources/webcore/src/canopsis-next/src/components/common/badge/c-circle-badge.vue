<template>
  <span
    :class="{
      'c-circle-badge--outlined': outlined,
      'c-circle-badge--small': small,
    }"
    :style="badgeStyle"
    class="c-circle-badge"
  >
    <v-progress-circular
      v-if="pending"
      size="12"
      width="1"
      indeterminate
    />
    <slot v-else />
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
    small: {
      type: Boolean,
      default: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
  },

  computed: {
    badgeStyle() {
      const color = this.color === 'white' ? 'white' : `rgb(var(--v-theme-${this.color}))`;

      return {
        '--c-circle-badge-color': color,
      };
    },
  },
};
</script>

<style lang="scss" scoped>
.c-circle-badge {
  --c-circle-badge-size: 22px;
  --c-circle-badge-font-size: calc(var(--c-circle-badge-size) * 0.6);

  align-items: center;
  background: var(--c-circle-badge-color);
  border: 1px solid var(--c-circle-badge-color);
  border-radius: calc(var(--c-circle-badge-size) / 2);
  color: rgb(var(--v-theme-on-primary));
  display: inline-flex;
  font-size: var(--c-circle-badge-font-size);
  height: var(--c-circle-badge-size);
  justify-content: center;
  min-width: var(--c-circle-badge-size);
  padding: 4px;

  &--outlined {
    background: transparent;
    color: var(--c-circle-badge-color);
  }

  &--small {
    --c-circle-badge-size: 20px;
    --c-circle-badge-font-size: calc(var(--c-circle-badge-size) * 0.6);
  }
}
</style>
