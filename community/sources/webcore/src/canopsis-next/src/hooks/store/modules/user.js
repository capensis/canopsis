import { useStoreModuleHooks } from '@/hooks/store';
import { usePopups } from '@/hooks/popups';
import { useI18n } from '@/hooks/i18n';

/**
 * Hook for accessing the user Vuex store module
 * Creates hooks for accessing user module's getters and actions
 *
 * @returns {Object} An object containing store module utilities
 * @property {Object} store - The Vuex store instance
 * @property {Object} module - The user Vuex module
 * @property {Function} useGetters - Function to access user module getters
 * @property {Function} useActions - Function to access user module actions
 * @example
 * const { useGetters, useActions } = useUserStoreModule();
 * const { users } = useGetters(['users']);
 * const { fetchUsers } = useActions(['fetchUsers']);
 */
const useUserStoreModule = () => useStoreModuleHooks('user');

/**
 * Hook for managing user-related operations and state
 *
 * @returns {Object} An object containing user-related getters, actions, and methods
 * @property {Array} users - List of users
 * @property {boolean} usersPending - Loading state of users
 * @property {Object} usersMeta - Metadata for users
 * @property {Function} fetchUsersList - Fetches the list of users
 * @property {Function} fetchUsersListWithPreviousParams - Fetches users list with previous parameters
 * @property {Function} createUser - Creates a new user
 * @property {Function} updateUser - Updates an existing user
 * @property {Function} updateCurrentUser - Updates the current user
 * @property {Function} removeUser - Removes a user
 * @property {Function} createUserWithPopup - Creates a user with success popup
 * @property {Function} updateUserWithPopup - Updates a user with success popup
 * @property {Function} removeUserWithPopup - Removes a user with success/error popup
 */
export const useUser = () => {
  const { t } = useI18n();
  const popups = usePopups();

  const { useGetters, useActions } = useUserStoreModule();

  const getters = useGetters({
    users: 'items',
    usersPending: 'pending',
    usersMeta: 'meta',
  });

  const actions = useActions({
    fetchUsersList: 'fetchList',
    fetchUsersListWithPreviousParams: 'fetchListWithPreviousParams',
    createUser: 'create',
    updateUser: 'update',
    updateCurrentUser: 'updateCurrentUser',
    removeUser: 'remove',
  });

  /**
   * Creates a new user and shows a success popup message
   * @async
   * @param {Object} params - The parameters object
   * @param {Object} params.data - The user data to create
   * @returns {Promise<void>}
   */
  const createUserWithPopup = async ({ data }) => {
    await actions.createUser({ data });

    popups.success({ text: t('success.default') });
  };

  /**
   * Updates an existing user and shows a success popup message
   * @async
   * @param {Object} params - The parameters object
   * @param {string|number} params.id - The ID of the user to update
   * @param {Object} params.data - The user data to update
   * @returns {Promise<void>}
   */
  const updateUserWithPopup = async ({ id, data }) => {
    await actions.updateUser({ id, data });

    popups.success({ text: t('success.default') });
  };

  /**
   * Removes a user and shows a success/error popup message
   * @async
   * @param {Object} params - The parameters object
   * @param {string|number} params.id - The ID of the user to remove
   * @returns {Promise<void>}
   * @throws {Error} When user removal fails
   */
  const removeUserWithPopup = async ({ id }) => {
    try {
      await actions.removeUser({ id });
      popups.success({ text: t('success.default') });
    } catch (err) {
      console.error(err);

      popups.error({ text: t('errors.default') });
    }
  };

  return {
    ...getters,
    ...actions,

    createUserWithPopup,
    updateUserWithPopup,
    removeUserWithPopup,
  };
};
