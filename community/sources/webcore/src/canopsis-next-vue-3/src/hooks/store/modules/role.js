import { storeToRefs } from 'pinia';

import { useRoleStore } from '@/stores/modules/role';
import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';

/**
 * Creates hooks for accessing role-related store operations.
 *
 * @returns {Object} Role refs and actions.
 */
export const useRole = () => {
  const { t } = useI18n();
  const popups = usePopups();
  const roleStore = useRoleStore();
  const {
    items,
    meta,
    pending,
    fetchingParams,
  } = storeToRefs(roleStore);

  /**
   * Creates a new role and shows a success popup message.
   *
   * @param {Object} params
   * @param {Object} params.data
   * @returns {Promise<void>}
   */
  const createRoleWithPopup = async ({ data }) => {
    await roleStore.create({ data });

    popups.success({ text: t('success.default') });
  };

  /**
   * Updates an existing role and shows a success popup message.
   *
   * @param {Object} params
   * @param {string|number} params.id
   * @param {Object} params.data
   * @returns {Promise<void>}
   */
  const updateRoleWithPopup = async ({ id, data }) => {
    await roleStore.update({ id, data });

    popups.success({ text: t('success.default') });
  };

  /**
   * Removes a role and shows a success/error popup message.
   *
   * @param {Object} params
   * @param {string|number} params.id
   * @returns {Promise<void>}
   */
  const removeRoleWithPopup = async ({ id }) => {
    try {
      await roleStore.remove({ id });
      popups.success({ text: t('success.default') });
    } catch (err) {
      console.error(err);

      popups.error({ text: t('errors.default') });
    }
  };

  return {
    roles: items,
    rolesMeta: meta,
    rolesPending: pending,
    rolesFetchingParams: fetchingParams,
    getRoleById: roleStore.getItemById,

    fetchRolesListWithoutStore: roleStore.fetchListWithoutStore,
    fetchRolesList: roleStore.fetchList,
    fetchRolesListWithPreviousParams: roleStore.fetchListWithPreviousParams,
    removeRole: roleStore.remove,
    createRole: roleStore.create,
    updateRole: roleStore.update,
    fetchRoleTemplatesListWithoutStore: roleStore.fetchTemplatesListWithoutStore,
    bulkUpdateRolePermissions: roleStore.bulkUpdatePermissions,

    createRoleWithPopup,
    updateRoleWithPopup,
    removeRoleWithPopup,
  };
};
