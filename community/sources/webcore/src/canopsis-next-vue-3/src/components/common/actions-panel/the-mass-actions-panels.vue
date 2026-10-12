<template>
  <div :data-teleport-target="TELEPORTS_NAMES.massActionsPanel" :class="targetClass" />
</template>

<script>
import { computed } from 'vue';

import { TELEPORTS_NAMES } from '@/constants';

export default {
  props: {
    fixed: {
      type: Boolean,
      default: false,
    },
  },

  setup(props) {
    const targetClass = computed(() => ({
      'the-mass-actions-panels__teleport-target': true,
      'the-mass-actions-panels__teleport-target--fixed': props.fixed,
    }));

    return {
      TELEPORTS_NAMES,
      targetClass,
    };
  },
};
</script>

<style lang="scss">
.the-mass-actions-panels__teleport-target {
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
