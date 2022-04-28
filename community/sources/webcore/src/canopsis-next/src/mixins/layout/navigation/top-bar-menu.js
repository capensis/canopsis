import { sortBy } from 'lodash';

import { authMixin } from '@/mixins/auth';
import { entitiesInfoMixin } from '@/mixins/entities/info';

export const layoutNavigationTopBarMenuMixin = {
  mixins: [authMixin, entitiesInfoMixin],
  methods: {
    prepareLinks(links) {
      const { permissionsWithDefaultType = [] } = this;
      const preparedLinks = links
        .filter(({ permission }) => {
          if (this.checkAppInfoAccessByPermission(permission)) {
            return permissionsWithDefaultType.includes(permission)
              ? this.checkAccess(permission)
              : this.checkReadAccess(permission);
          }

          return false;
        })
        .map(link => ({ ...link, title: this.$t(`pageHeaders.${link.permission}.title`) }));

      return sortBy(preparedLinks, 'title');
    },
  },
};
