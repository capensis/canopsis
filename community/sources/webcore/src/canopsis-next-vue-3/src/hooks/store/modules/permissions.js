import { usePermissionStore } from '@/stores/modules/permission';

/**
 * Hook for accessing permission-related actions.
 *
 * @returns {{ fetchPermissionsListWithoutStore: Function }}
 */
export const usePermissions = () => {
  const permissionStore = usePermissionStore();

  return {
    fetchPermissionsListWithoutStore: permissionStore.fetchListWithoutStore,
  };
};
