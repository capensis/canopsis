import { sortBy } from 'lodash';

import { authMixin } from '@/mixins/auth';
import { entitiesInfoMixin } from '@/mixins/entities/info';

export const layoutNavigationTopBarMenuMixin = {
  mixins: [authMixin, entitiesInfoMixin],
  methods: {
    filterLinks(links) {
      const { permissionsWithDefaultType = [] } = this;

      return links
        .filter(({ permission }) => {
          if (!permission) {
            return true;
          }

          if (this.checkAppInfoAccessByPermission(permission)) {
            return permissionsWithDefaultType.includes(permission)
              ? this.checkAccess(permission)
              : this.checkReadAccess(permission);
          }

          return false;
        });
    },

    prepareLinks(links) {
      const preparedLinks = this.filterLinks(links)
        .map(link => ({ ...link, title: this.$t(`pageHeaders.${link.permission}.title`) }));

      return sortBy(preparedLinks, 'title');
    },
  },
};
