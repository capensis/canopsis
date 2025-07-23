import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the activeView Vuex store module.
 * Provides access to getters and actions within the 'activeView' namespace.
 *
 * @returns {Object} An object containing:
 *   - store: The Vuex store instance
 *   - module: The activeView module instance
 *   - useGetters: Function to access module getters
 *   - useActions: Function to access module actions
 */
const useActiveViewStoreModule = () => useStoreModuleHooks('activeView');

/**
 * Hook for managing active view-related operations and state
 * Replaces the activeViewMixin functionality with Composition API
 *
 * @returns {Object} An object containing active view getters and actions
 * @property {Object} view - The active view item
 * @property {boolean} pending - Loading state of the view
 * @property {boolean} editing - Whether the view is in editing mode
 * @property {string} mode - The current view mode
 * @property {boolean} editingProcess - Whether an editing process is ongoing
 * @property {Function} toggleEditing - Toggles the editing state
 * @property {Function} registerEditingOffHandler - Registers a handler for editing off
 * @property {Function} unregisterEditingOffHandler - Unregisters the editing off handler
 * @property {Function} fetchActiveView - Fetches the active view data
 * @property {Function} clearActiveView - Clears the active view data
 * @property {Function} setActiveViewMode - Sets the active view mode
 */
export const useActiveView = () => {
  const { useGetters, useActions } = useActiveViewStoreModule();

  const getters = useGetters({
    view: 'item',
    pending: 'pending',
    editing: 'editing',
    screenMode: 'screenMode',
    isKioskScreenMode: 'isKioskScreenMode',
    editingProcess: 'editingProcess',
  });

  const actions = useActions({
    toggleEditing: 'toggleEditing',
    registerEditingOffHandler: 'registerEditingOffHandler',
    unregisterEditingOffHandler: 'unregisterEditingOffHandler',
    fetchActiveView: 'fetch',
    clearActiveView: 'clear',
    setActiveViewScreenMode: 'setScreenMode',
  });

  return {
    ...getters,
    ...actions,
  };
};
