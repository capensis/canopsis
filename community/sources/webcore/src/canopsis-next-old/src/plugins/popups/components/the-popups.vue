<template>
  <div
    v-if="popups.length"
    class="popups"
  >
    <popup-item
      v-for="popup in popups"
      v-bind="popup"
      :key="popup.id"
    />
  </div>
</template>

<script>
import { computed } from 'vue';

import { useStore } from '@/hooks/store';
import { usePopups } from '@/hooks/popups';

import PopupItem from './popup-item.vue';

/**
 * Wrapper for the popups
 */
export default {
  components: { PopupItem },
  setup() {
    const store = useStore();
    const popupsModule = usePopups();

    const popups = computed(() => store.getters[`${popupsModule.moduleName}/popups`]);

    return {
      popups,
    };
  },
};
</script>

<style lang="scss">
  .popups {
    position: fixed;
    z-index: 1000;
    right: 2rem;
    top: 75px;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
  }
</style>
