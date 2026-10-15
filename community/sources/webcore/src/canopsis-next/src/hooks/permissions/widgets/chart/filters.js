import { computed } from 'vue';

import { USER_PERMISSIONS } from '@/constants';

import { useCurrentUserPermissions } from '@/hooks/auth';

/**
 * Hook that checks current user permissions for chart widget filter settings.
 *
 * @param {'barChart' | 'lineChart' | 'pieChart' | 'numbers'} permissionKey
 * Key from `USER_PERMISSIONS.business` identifying the chart widget type.
 *
 * @returns {Object} Chart widget filter permission state.
 * @property {import('vue').ComputedRef<boolean>} hasAccessToFilter
 * Whether the user can access the widget filter setting.
 * @property {import('vue').ComputedRef<boolean>} hasAccessToUserFilter
 * Whether the user can access the user filter setting.
 */
const useChartWidgetFiltersPermissions = (permissionKey) => {
  const { checkAccess } = useCurrentUserPermissions();

  const hasAccessToFilter = computed(() =>
    checkAccess(USER_PERMISSIONS.business[permissionKey].actions.filter),
  );

  const hasAccessToUserFilter = computed(() =>
    checkAccess(USER_PERMISSIONS.business[permissionKey].actions.userFilter),
  );

  return {
    hasAccessToFilter,
    hasAccessToUserFilter,
  };
};

/** @returns {ReturnType<typeof useChartWidgetFiltersPermissions>} Bar chart widget filter permissions. */
export const useBarChartFiltersPermissions = () => useChartWidgetFiltersPermissions('barChart');

/** @returns {ReturnType<typeof useChartWidgetFiltersPermissions>} Line chart widget filter permissions. */
export const useLineChartFiltersPermissions = () => useChartWidgetFiltersPermissions('lineChart');

/** @returns {ReturnType<typeof useChartWidgetFiltersPermissions>} Pie chart widget filter permissions. */
export const usePieChartFiltersPermissions = () => useChartWidgetFiltersPermissions('pieChart');

/** @returns {ReturnType<typeof useChartWidgetFiltersPermissions>} Numbers chart widget filter permissions. */
export const useNumbersChartFiltersPermissions = () => useChartWidgetFiltersPermissions('numbers');
