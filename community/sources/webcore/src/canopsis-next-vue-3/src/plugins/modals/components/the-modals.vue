<script setup>
import { storeToRefs } from 'pinia';
import { computed, watch } from 'vue';
import { useRoute } from 'vue-router';

import { useModalsStore } from '../stores';

const route = useRoute();
const modalsStore = useModalsStore();

const { modals } = storeToRefs(modalsStore);

const zIndex = computed(() => modals.value.length * 1010);

/**
 * Hide all non-minimized modals.
 */
const hideOpenedModals = () => {
  modals.value.forEach((modal) => {
    if (!modal.minimized) {
      modalsStore.hide({ id: modal.id });
    }
  });
};

watch(
  () => route.path,
  (newPath, prevPath) => {
    if (newPath === prevPath || !modals.value.length) {
      return;
    }

    hideOpenedModals();
  },
);
</script>

<template>
  <div
    :class="{ 'modals-wrapper--active': modals.length > 0 }"
    :style="{ zIndex }"
    class="modals-wrapper"
  >
    <modal-base v-for="modal in modals" :key="modal.id" :modal="modal" />

    <div class="pbehavior-calendar-event-menu" />
  </div>
</template>

<style lang="scss" scoped>
$minimizedDialogMaxWidth: 360px;

.modals-wrapper {
  display: none;

  .pbehavior-calendar-event-menu {
    position: absolute;
    inset: 0;
    display: flex;
    align-items: center;
  }

  &--active {
    position: fixed;
    display: flex;
    bottom: 0;
    width: 100%;
    height: 100%;
    flex-wrap: wrap-reverse;
    pointer-events: none;
    justify-content: center;
    align-items: flex-end;
    align-content: flex-start;
    z-index: 300;

    :deep(.v-overlay__content) {
      pointer-events: auto;
    }
  }
}
</style>
