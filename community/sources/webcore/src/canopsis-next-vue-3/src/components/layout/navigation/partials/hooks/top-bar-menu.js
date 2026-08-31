import { isArray, sortBy } from 'lodash';
import { computed, unref } from 'vue';

import { groupedPermissionToPermission } from '@/helpers/permission';

import { useCurrentUserPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useInfo } from '@/hooks/store/modules/info';

export const useTopBarMenu = ({ permissionsWithDefaultType = [], withoutSort = true } = {}) => {
  const { currentUserPermissionsById, checkAccess, checkReadAccess } = useCurrentUserPermissions();
  const { checkAppInfoAccessByPermission } = useInfo();
  const { t, te, tm } = useI18n();

  const hasLoadedPermissions = computed(
    () => Object.keys(currentUserPermissionsById.value ?? {}).length > 0,
  );

  const checkPermissionAccess = (permission) => {
    if (!hasLoadedPermissions.value) {
      return true;
    }

    return unref(permissionsWithDefaultType).includes(permission)
      ? checkAccess(permission)
      : checkReadAccess(permission);
  };

  const checkWholeAccess = (permission) => {
    if (!checkAppInfoAccessByPermission(permission)) {
      return false;
    }

    if (isArray(permission)) {
      return permission.some(checkPermissionAccess);
    }

    return checkPermissionAccess(permission);
  };

  const filterLinks = (links = []) => links.filter(({ permission }) => {
      if (!permission) {
        return true;
      }

      return checkWholeAccess(permission);
    });

  const getPageHeader = (permission) => {
    if (!permission) {
      return {};
    }

    return tm(`pageHeaders.${permission}`) || {};
  };

  const getPermissionTitle = (permission) => {
    if (!permission) {
      return '';
    }

    const key = `permission.title.${permission}`;

    return te(key) ? t(key) : '';
  };

  const prepareLinks = (links = []) => {
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

        if (link.divider || link.header) {
          return link;
        }

        const permissionName = isArray(link.permission)
          ? groupedPermissionToPermission(link.permission)
          : link.permission;
        const { topbarTitle, title } = getPageHeader(permissionName);

        return {
          ...link,
          title:
            link.title ||
            topbarTitle ||
            title ||
            getPermissionTitle(permissionName) ||
            link.route?.name ||
            '',
        };
      })
      .filter(Boolean)
      .filter((link, index, currentLinks) => {
        if (link.header) {
          const nextLink = currentLinks[index + 1];

          return !nextLink?.divider;
        }

        if (link.divider) {
          const prevLink = currentLinks[index - 1];

          return !prevLink?.header && !prevLink?.divider;
        }

        return true;
      });

    return unref(withoutSort) ? preparedLinks : sortBy(preparedLinks, 'title');
  };

  return {
    checkWholeAccess,
    filterLinks,
    prepareLinks,
  };
};
