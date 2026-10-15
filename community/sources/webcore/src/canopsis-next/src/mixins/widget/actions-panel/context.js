import { BUSINESS_USER_PERMISSIONS_ACTIONS_MAP, ENTITY_TYPES } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { entitiesContextEntityMixin } from '@/mixins/entities/context-entity';
import { entitiesServiceMixin } from '@/mixins/entities/service';

export const widgetActionsPanelContextMixin = {
  mixins: [authMixin, entitiesServiceMixin, entitiesContextEntityMixin],
  methods: {
    actionsAccessFilterHandler({ type }) {
      const permission = BUSINESS_USER_PERMISSIONS_ACTIONS_MAP.context[type];

      if (!permission) {
        return true;
      }

      return this.checkAccess(permission);
    },

    /**
     * Remove context entity or service by item
     *
     * @param {Entity | Service | Object} item
     * @returns {Promise}
     */
    removeContextEntityOrService(item) {
      if (item.type === ENTITY_TYPES.service) {
        return this.removeService({ id: item._id });
      }

      return this.removeContextEntity({ id: item._id });
    },
  },
};
