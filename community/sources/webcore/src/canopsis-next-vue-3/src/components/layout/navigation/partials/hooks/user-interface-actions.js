import { useModals } from '@/hooks/modals';
import { MODALS } from '@/constants';

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
