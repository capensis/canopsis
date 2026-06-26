import { ref } from 'vue';
import { defineStore } from 'pinia';

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
    pending.value = false;
    editing.value = false;
    screenMode.value = VIEW_SCREEN_MODES.default;
    editingOffHandlers.value = [];
    periodicRefreshPaused.value = false;
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

    fetch,
    clear,
  };
});
