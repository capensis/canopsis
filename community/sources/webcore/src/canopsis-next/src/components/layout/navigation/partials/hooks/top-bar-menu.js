import { sortBy, isArray } from 'lodash';
import { unref } from 'vue';

import { groupedPermissionToPermission } from '@/helpers/permission';

import { useCurrentUserPermissions } from '@/hooks/auth';
import { useInfo } from '@/hooks/store/modules/info';
import { useI18n } from '@/hooks/i18n';

/**
 * Hook for top bar menu functionality
 * Provides methods to filter and prepare menu links based on user permissions
 *
 * @param {Object} options - Configuration options
 * @param {Array} [options.permissionsWithDefaultType=[]] - List of permissions that use default type checking
 * @param {boolean} [options.withoutSort=true] - Whether to sort links by title
 * @returns {Object} An object containing methods to prepare menu links
 */
export const useTopBarMenu = ({ permissionsWithDefaultType = [], withoutSort = true } = {}) => {
  const { checkAccess, checkReadAccess } = useCurrentUserPermissions();
  const { checkAppInfoAccessByPermission } = useInfo();
  const { t } = useI18n();

  /**
   * Check if user has whole access to a permission
   *
   * @param {string|Array} permission - Permission or array of permissions to check
   * @returns {boolean} Whether the user has access
   */
  const checkWholeAccess = (permission) => {
    if (checkAppInfoAccessByPermission(permission)) {
      return unref(permissionsWithDefaultType).includes(permission)
        ? checkAccess(permission)
        : checkReadAccess(permission);
    }

    return false;
  };

  /**
   * Filter links based on user permissions
   *
   * @param {Array} [links=[]] - Array of link objects
   * @returns {Array} Filtered array of links
   */
  const filterLinks = (links = []) => links.filter(({ permission }) => {
    if (!permission) {
      return true;
    }

    const arrayPermissions = isArray(permission) ? permission : [permission];

    return arrayPermissions.some(checkWholeAccess);
  });

  /**
   * Prepare links by filtering and adding translated titles
   *
   * @param {Array} links - Array of link objects
   * @returns {Array} Prepared and sorted array of links
   */
  const prepareLinks = (links) => {
    const preparedLinks = filterLinks(links)
      .map((link) => {
        if (link.links) {
          const nestedLinks = prepareLinks(link.links);

          if (!nestedLinks.length) {
            return null;
          }

          return {
            ...link,
            links: nestedLinks,
          };
        }

        const permissionName = isArray(link.permission)
          ? groupedPermissionToPermission(link.permission)
          : link.permission;

        const { topbarTitle, title } = t(`pageHeaders.${permissionName}`);

        return {
          ...link,
          title: link.title || topbarTitle || title,
        };
      })
      .filter(Boolean);

    return unref(withoutSort) ? preparedLinks : sortBy(preparedLinks, 'title');
  };

  return {
    checkWholeAccess,
    filterLinks,
    prepareLinks,
  };
};
