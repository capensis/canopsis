import { storeToRefs } from 'pinia';

import { useActiveViewStore } from '@/stores/modules/active-view';

/**
 * Hook for managing active view-related operations and state.
 *
 * @returns {Object} Active view refs and actions.
 */
export const useActiveView = () => {
  const activeViewStore = useActiveViewStore();
  const {
    item,
    pending,
    screenMode,
    editing,
    editingProcess,
    periodicRefreshPaused,
    isKioskScreenMode,
  } = storeToRefs(activeViewStore);

  return {
    view: item,
    activeViewPending: pending,
    activeViewScreenMode: screenMode,
    activeViewEditing: editing,
    activeViewEditingProcess: editingProcess,
    activeViewPeriodicRefreshPaused: periodicRefreshPaused,
    activeViewIsKioskScreenMode: isKioskScreenMode,

    toggleEditing: activeViewStore.toggleEditing,
    registerEditingOffHandler: activeViewStore.registerEditingOffHandler,
    unregisterEditingOffHandler: activeViewStore.unregisterEditingOffHandler,
    fetchActiveView: activeViewStore.fetch,
    clearActiveView: activeViewStore.clear,
    resumePeriodicRefresh: activeViewStore.resumePeriodicRefresh,
    pausePeriodicRefresh: activeViewStore.pausePeriodicRefresh,
    setActiveViewScreenMode: activeViewStore.setScreenMode,
  };
};
