import { USERS_RIGHTS_MASKS, USERS_RIGHTS_TYPES } from '@/constants';
import { generateRight } from '@/helpers/entities';

import rightEntitiesBaseMixin from '@/mixins/rights/entities/base';

export default {
  mixins: [rightEntitiesBaseMixin],
  methods: {
    async createRightByViewId(viewId) {
      const checksum = USERS_RIGHTS_MASKS.read + USERS_RIGHTS_MASKS.update + USERS_RIGHTS_MASKS.delete;
      const right = {
        ...generateRight(),

        _id: viewId,
        type: USERS_RIGHTS_TYPES.rw,
        desc: `Rights on view: ${viewId}`,
      };

      return this.createRightAndAddIntoRole(right, checksum);
    },

    async removeRightByViewId(viewId) {
      return this.removeRightAndRemoveRightFromRoleById(viewId);
    },
  },
};
