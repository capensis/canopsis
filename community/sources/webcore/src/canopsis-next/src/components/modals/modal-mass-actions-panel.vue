<template>
  <portal-target
    :name="name"
    class="modal-mass-actions-panel__portal-target"
  />
</template>

<script>
import { computed } from 'vue';

import { PORTALS_NAMES } from '@/constants';

export default {
  props: {
    id: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const name = computed(() => `${PORTALS_NAMES.massActionsPanel}-${props.id}`);

    return {
      name,
    };
  },
};
</script>

<style lang="scss" scoped>
.modal-mass-actions-panel__portal-target {
  position: absolute;
  right: 0;
  top: 0;
  width: 100%;
  z-index: 4;
  pointer-events: none;
}
</style>
