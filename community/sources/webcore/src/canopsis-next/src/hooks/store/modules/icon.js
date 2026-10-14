import { storeToRefs } from 'pinia';

import { useIconStore } from '@/stores/modules/icon';

/**
 * Icon store facade.
 *
 * @returns {Object}
 */
export const useIcon = () => {
  const iconStore = useIconStore();
  const { registeredIconsById } = storeToRefs(iconStore);

  return {
    registeredIconsById,

    fetchIconsListWithoutStore: iconStore.fetchListWithoutStore,
    fetchIconWithoutStore: iconStore.fetchItemWithoutStore,
    createIcon: iconStore.create,
    updateIcon: iconStore.update,
    removeIcon: iconStore.remove,
    addRegisteredIcon: iconStore.addRegisteredIcon,
    removeRegisteredIcon: iconStore.removeRegisteredIcon,
    setRegisteredIcons: iconStore.setRegisteredIcons,
  };
};
