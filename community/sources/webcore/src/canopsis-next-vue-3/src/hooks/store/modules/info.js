import { computed } from 'vue';
import { storeToRefs } from 'pinia';
import { isArray, isMatch } from 'lodash';

import { DEFAULT_APP_TITLE } from '@/config';
import { CANOPSIS_EDITION, USER_PERMISSIONS_TO_PAGES_RULES } from '@/constants';

import { compile } from '@/helpers/handlebars';
import { sanitizeHtml } from '@/helpers/html';
import { groupedPermissionToPermission } from '@/helpers/permission';

import { useAppInfoStore } from '@/stores/modules/app-info';

/**
 * App info store facade.
 *
 * Keeps Pinia access isolated from components and exposes app-info state,
 * derived flags and actions through a migration-friendly hook API.
 *
 * @returns {Object}
 */
export const useInfo = () => {
  const appInfoStore = useAppInfoStore();
  const refs = storeToRefs(appInfoStore);

  const appInfoPending = refs.pending;
  const isProVersion = computed(() => refs.edition.value === CANOPSIS_EDITION.pro);

  const checkAppInfoAccessByPermission = (permission) => {
    const preparedPermission = isArray(permission)
      ? groupedPermissionToPermission(permission)
      : permission;
    const permissionAppInfoRules = USER_PERMISSIONS_TO_PAGES_RULES[preparedPermission];

    if (!permissionAppInfoRules) {
      return true;
    }

    return isMatch({ edition: refs.edition.value }, permissionAppInfoRules);
  };

  const setTitle = async () => {
    document.title = refs.appTitle.value
      ? await compile(
          sanitizeHtml(refs.appTitle.value, {
            ALLOWED_TAGS: [],
            ALLOWED_ATTR: [],
          }),
        )
      : DEFAULT_APP_TITLE;
  };

  return {
    ...refs,
    appInfoPending,
    isProVersion,

    fetchAppInfo: appInfoStore.fetchAppInfo,
    updateUserInterface: appInfoStore.updateUserInterface,
    updateMaintenanceMode: appInfoStore.updateMaintenanceMode,
    checkAppInfoAccessByPermission,
    setTitle,
  };
};
