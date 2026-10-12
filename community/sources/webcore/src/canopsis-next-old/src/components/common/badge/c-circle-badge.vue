<template>
  <span :class="classes" class="c-circle-badge">
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
import { computed } from 'vue';

export default {
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
  setup(props) {
    const classes = computed(() => ({
      'c-circle-badge--small': props.small,
      'c-circle-badge--outlined': props.outlined,
      [`${props.color}--text`]: props.outlined,
      [props.color]: true,
    }));

    return {
      classes,
    };
  },
};
</script>

<style lang="scss">
.c-circle-badge {
  --c-circle-badge-size: 22px;
  --c-circle-badge-font-size: calc(var(--c-circle-badge-size) * 0.6);

  color: white;

  padding: 4px;

  display: inline-flex;
  align-items: center;
  justify-content: center;

  font-size: var(--c-circle-badge-font-size);

  height: var(--c-circle-badge-size);
  min-width: var(--c-circle-badge-size);
  border-radius: calc(var(--c-circle-badge-size) / 2);

  &--small {
    --c-circle-badge-size: 20px;
    --c-circle-badge-font-size: calc(var(--c-circle-badge-size) * 0.6);
  }

  .v-application &--outlined.c-circle-badge {
    background-color: transparent !important;
    border: 1px solid currentColor;
  }
}
</style>
