import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { VIEW_SCREEN_MODES } from '@/constants';

import { useViewStore } from './view';

export const useActiveViewStore = defineStore('activeView', () => {
  const id = ref(null);
  const item = ref(null);
  const pending = ref(false);
  const editing = ref(false);
  const editingProcess = ref(false);
  const screenMode = ref(VIEW_SCREEN_MODES.default);
  const editingOffHandlers = ref([]);
  const periodicRefreshPaused = ref(false);

  const isKioskScreenMode = computed(() =>
    [VIEW_SCREEN_MODES.kiosk, VIEW_SCREEN_MODES.kioskFullscreen].includes(screenMode.value),
  );

  const registerEditingOffHandler = (handler) => {
    editingOffHandlers.value.push(handler);
  };

  const unregisterEditingOffHandler = (handler) => {
    editingOffHandlers.value = editingOffHandlers.value.filter(
      editingHandler => editingHandler !== handler,
    );
  };

  const toggleEditing = async () => {
    try {
      editingProcess.value = true;

      if (editingOffHandlers.value.length && editing.value) {
        await Promise.all(editingOffHandlers.value.map(handler => handler()));
      }

      editing.value = !editing.value;
      editingProcess.value = false;
    } catch (err) {
      console.error(err);
    }
  };

  const fetch = async ({ id: viewId = id.value } = {}) => {
    try {
      if (!viewId) {
        throw new Error('Active view id is empty');
      }

      pending.value = true;
      id.value = viewId;

      const viewStore = useViewStore();
      const view = await viewStore.fetchView({ id: viewId });

      item.value = view ?? item.value;

      return item.value;
    } catch (err) {
      pending.value = false;
      id.value = null;

      console.error(err);

      throw err;
    } finally {
      pending.value = false;
    }
  };

  const clear = () => {
    id.value = null;
    item.value = null;
    pending.value = false;
    editing.value = false;
    screenMode.value = VIEW_SCREEN_MODES.default;
    editingOffHandlers.value = [];
    periodicRefreshPaused.value = false;
  };

  const resumePeriodicRefresh = () => {
    periodicRefreshPaused.value = false;
  };

  const pausePeriodicRefresh = () => {
    periodicRefreshPaused.value = true;
  };

  const setScreenMode = (newScreenMode) => {
    if (screenMode.value === newScreenMode) {
      return;
    }

    screenMode.value = newScreenMode;
  };

  return {
    id,
    item,
    pending,
    editing,
    editingProcess,
    screenMode,
    editingOffHandlers,
    periodicRefreshPaused,
    isKioskScreenMode,

    registerEditingOffHandler,
    unregisterEditingOffHandler,
    toggleEditing,
    fetch,
    clear,
    resumePeriodicRefresh,
    pausePeriodicRefresh,
    setScreenMode,
  };
});
