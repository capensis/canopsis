import { useAuth } from '@/hooks/auth';
import { useUser } from '@/hooks/store/modules/user';

/**
 * Tour completion helpers: updates `currentUser.ui_tours` via the user API and refreshes auth state.
 *
 * @returns {{ currentUser: import('vue').Ref<Object>, finishTourByName: Function }}
 */
export const useTourBase = () => {
  const { currentUser } = useAuth();
  const { updateCurrentUserTours } = useUser();

  /**
   * Marks a named page tour as finished for the current user and refetches profile.
   *
   * @param {string} tourName - Key stored under `ui_tours` (e.g. page header `name`).
   */
  const finishTourByName = tourName => updateCurrentUserTours({ data: { [tourName]: true } });

  return {
    currentUser,
    finishTourByName,
    updateCurrentUserTours,
  };
};
