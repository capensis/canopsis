import { storeToRefs } from 'pinia';

import { useViewStore } from '@/stores/modules/view';

/**
 * Creates hooks for accessing view group-related store operations.
 *
 * @returns {Object} View group refs and actions.
 */
export const useView = () => {
  const viewStore = useViewStore();
  const { items, meta, pending } = storeToRefs(viewStore);

  return {
    viewGroups: items,
    viewGroupsMeta: meta,
    viewGroupsPending: pending,

    getViewGroupById: viewStore.getGroupById,
    getViewById: viewStore.getViewById,
    getViewTabById: viewStore.getViewTabById,

    fetchViewGroupsList: viewStore.fetchList,
    fetchViewGroupsListWithoutStore: viewStore.fetchListWithoutStore,
    createViewGroup: viewStore.create,
    updateViewGroup: viewStore.update,
    removeViewGroup: viewStore.remove,
    createPrivateViewGroup: viewStore.createPrivateGroup,
    updatePrivateViewGroup: viewStore.updatePrivateGroup,
    removePrivateViewGroup: viewStore.removePrivateGroup,
    createView: viewStore.createView,
    updateView: viewStore.updateView,
    updateViewWithoutStore: viewStore.updateWithoutStoreView,
    updateViewPositions: viewStore.updateViewPositions,
    removeView: viewStore.removeView,
    copyView: viewStore.copyView,
    exportViewWithoutStore: viewStore.exportViewWithoutStore,
    importViewWithoutStore: viewStore.importViewWithoutStore,
    fetchViewTab: viewStore.fetchViewTab,
    createViewTab: viewStore.createViewTab,
    updateViewTab: viewStore.updateViewTab,
    copyViewTab: viewStore.copyViewTab,
    removeViewTab: viewStore.removeViewTab,
    updateViewTabPositions: viewStore.updateViewTabPositions,
    fetchView: viewStore.fetchView,
  };
};
