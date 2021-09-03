import { sortBy } from 'lodash';

import { authMixin } from '@/mixins/auth';
import entitiesInfoMixin from '@/mixins/entities/info';

export const layoutNavigationTopBarMenuMixin = {
  mixins: [authMixin, entitiesInfoMixin],
  methods: {
    prepareLinks(links) {
      const { permissionsWithDefaultType = [] } = this;
      const preparedLinks = links
        .filter(({ permission }) => (permissionsWithDefaultType.includes(permission)
          ? this.checkAccess(permission)
          : this.checkAppInfoAccessByPermission(permission) && this.checkReadAccess(permission)))
        .map(link => ({ ...link, title: this.$t(`pageHeaders.${link.permission}.title`) }));

      return sortBy(preparedLinks, 'title');
    },
  },
};
