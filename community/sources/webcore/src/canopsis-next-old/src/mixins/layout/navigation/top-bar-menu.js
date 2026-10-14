import { sortBy, isArray } from 'lodash';

import { groupedPermissionToPermission } from '@/helpers/permission';

import { authMixin } from '@/mixins/auth';
import { entitiesInfoMixin } from '@/mixins/entities/info';

export const layoutNavigationTopBarMenuMixin = {
  mixins: [authMixin, entitiesInfoMixin],
  methods: {
    checkWholeAccess(permission) {
      if (this.checkAppInfoAccessByPermission(permission)) {
        return (this.permissionsWithDefaultType ?? []).includes(permission)
          ? this.checkAccess(permission)
          : this.checkReadAccess(permission);
      }

      return false;
    },

    filterLinks(links) {
      return links
        .filter(({ permission }) => {
          if (!permission) {
            return true;
          }

          const arrayPermissions = isArray(permission) ? permission : [permission];

          return arrayPermissions.some(this.checkWholeAccess);
        });
    },

    prepareLinks(links) {
      const preparedLinks = this.filterLinks(links)
        .map((link) => {
          const permissionName = isArray(link.permission)
            ? groupedPermissionToPermission(link.permission)
            : link.permission;

          return {
            ...link,

            title: this.$t(`pageHeaders.${permissionName}.title`),
          };
        });

      return sortBy(preparedLinks, 'title');
    },
  },
};
