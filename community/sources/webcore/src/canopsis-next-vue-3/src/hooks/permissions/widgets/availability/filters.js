import { computed } from 'vue';

import { USER_PERMISSIONS } from '@/constants';

import { useCurrentUserPermissions } from '@/hooks/auth';

/**
 * Hook that checks current user permissions for availability widget filter settings.
 *
 * @returns {Object} Availability widget filter permission state.
 * @property {import('vue').ComputedRef<boolean>} hasAccessToFilter
 * Whether the user can access the widget filter setting.
 * @property {import('vue').ComputedRef<boolean>} hasAccessToUserFilter
 * Whether the user can access the user filter setting.
 */
export const useAvailabilityFiltersPermissions = () => {
  const { checkAccess } = useCurrentUserPermissions();

  const hasAccessToFilter = computed(() => checkAccess(USER_PERMISSIONS.business.availability.actions.filter),
  );

  const hasAccessToUserFilter = computed(() => checkAccess(USER_PERMISSIONS.business.availability.actions.userFilter),
  );

  return {
    hasAccessToFilter,
    hasAccessToUserFilter,
  };
};

/**
 * Hook that checks current user permissions for availability widget interval settings.
 *
 * @returns {Object} Availability widget interval permission state.
 * @property {import('vue').ComputedRef<boolean>} hasAccessToInterval
 * Whether the user can access the widget interval setting.
 */
export const useAvailabilityIntervalPermissions = () => {
  const { checkAccess } = useCurrentUserPermissions();

  const hasAccessToInterval = computed(() => checkAccess(USER_PERMISSIONS.business.availability.actions.interval),
  );

  return {
    hasAccessToInterval,
  };
};

/**
 * Hook that checks current user permissions for availability widget CSV export.
 *
 * @returns {Object} Availability widget export permission state.
 * @property {import('vue').ComputedRef<boolean>} hasAccessToExportAsCsv
 * Whether the user can export availability data as CSV.
 */
export const useAvailabilityExportPermissions = () => {
  const { checkAccess } = useCurrentUserPermissions();

  const hasAccessToExportAsCsv = computed(() => checkAccess(USER_PERMISSIONS.business.availability.actions.exportAsCsv),
  );

  return {
    hasAccessToExportAsCsv,
  };
};
