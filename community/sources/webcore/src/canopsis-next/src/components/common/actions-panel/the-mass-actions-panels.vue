<template>
  <portal-target
    :name="$constants.PORTALS_NAMES.massActionsPanel"
    :class="targetClass"
    multiple
  />
</template>

<script>
import { computed } from 'vue';

export default {
  props: {
    fixed: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const targetClass = computed(() => ({
      'the-mass-actions-panels__portal-target': true,
      'the-mass-actions-panels__portal-target--fixed': props.fixed,
    }));

    return {
      targetClass,
    };
  },
};
</script>

<style lang="scss">
.the-mass-actions-panels__portal-target {
  --actions-grid-columns-count: 12;
  --actions-grid-column-gap: 20px;
  --actions-grid-padding: 0 20px;

  position: absolute;
  right: 0;
  top: 0;
  width: 100%;
  z-index: 3;
  padding: var(--actions-grid-padding);
  column-gap: var(--actions-grid-column-gap);
  display: grid;
  grid-template-columns: repeat(var(--actions-grid-columns-count), 1fr);
  pointer-events: none;

  &--fixed {
    position: sticky;
    top: 0;

    .c-mass-actions-panel {
      padding-top: 0;
    }
  }
}
</style>
