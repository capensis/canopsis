import { MODALS } from '@/constants';

import { useModals } from '@/hooks/modals';

/**
 * Hook for user interface settings modal actions
 *
 * @returns {Object} An object containing method to show user interface modal
 */
export const useUserInterfaceActions = () => {
  const modals = useModals();

  const showUserInterfaceModal = () => modals.show({ name: MODALS.userInterface });

  return {
    showUserInterfaceModal,
  };
};
