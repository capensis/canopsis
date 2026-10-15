import { computed } from 'vue';

import { USER_PERMISSIONS } from '@/constants';

import { useCurrentUserPermissions } from '@/hooks/auth';

/**
 * Hook that checks current user permissions for statistics widget filter settings.
 *
 * @param {'userStatistics' | 'alarmStatistics'} permissionKey
 * Key from `USER_PERMISSIONS.business` identifying the statistics widget type.
 *
 * @returns {Object} Statistics widget filter permission state.
 * @property {import('vue').ComputedRef<boolean>} hasAccessToFilter
 * Whether the user can access the widget filter setting.
 * @property {import('vue').ComputedRef<boolean>} hasAccessToUserFilter
 * Whether the user can access the user filter setting.
 */
const useStatisticsWidgetFiltersPermissions = (permissionKey) => {
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

/**
 * Hook that checks current user permissions for statistics widget interval settings.
 *
 * @param {'userStatistics' | 'alarmStatistics'} permissionKey
 * Key from `USER_PERMISSIONS.business` identifying the statistics widget type.
 *
 * @returns {Object} Statistics widget interval permission state.
 * @property {import('vue').ComputedRef<boolean>} hasAccessToInterval
 * Whether the user can access the widget interval setting.
 */
const useStatisticsWidgetIntervalPermissions = (permissionKey) => {
  const { checkAccess } = useCurrentUserPermissions();

  const hasAccessToInterval = computed(() =>
    checkAccess(USER_PERMISSIONS.business[permissionKey].actions.interval),
  );

  return {
    hasAccessToInterval,
  };
};

/** @returns {ReturnType<typeof useStatisticsWidgetFiltersPermissions>} User statistics filter permissions. */
export const useUserStatisticsFiltersPermissions = () =>
  useStatisticsWidgetFiltersPermissions('userStatistics');

/** @returns {ReturnType<typeof useStatisticsWidgetIntervalPermissions>} User statistics interval permissions. */
export const useUserStatisticsIntervalPermissions = () =>
  useStatisticsWidgetIntervalPermissions('userStatistics');

/** @returns {ReturnType<typeof useStatisticsWidgetFiltersPermissions>} Alarm statistics filter permissions. */
export const useAlarmStatisticsFiltersPermissions = () =>
  useStatisticsWidgetFiltersPermissions('alarmStatistics');

/** @returns {ReturnType<typeof useStatisticsWidgetIntervalPermissions>} Alarm statistics interval permissions. */
export const useAlarmStatisticsIntervalPermissions = () =>
  useStatisticsWidgetIntervalPermissions('alarmStatistics');
