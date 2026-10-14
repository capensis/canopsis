import { MODALS } from '@/constants';

import { useModals } from '@/hooks/modals';

/**
 * Provides actions for the user interface settings menu entry.
 *
 * @returns {{ showUserInterfaceModal: Function }}
 */
export const useUserInterfaceActions = () => {
  const modals = useModals();

  const showUserInterfaceModal = () =>
    modals.show({
      name: MODALS.userInterface,
    });

  return {
    showUserInterfaceModal,
  };
};
