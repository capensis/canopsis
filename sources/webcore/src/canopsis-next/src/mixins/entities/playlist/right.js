import { USERS_RIGHTS_MASKS } from '@/constants';
import { generateRight } from '@/helpers/entities';

import rightEntitiesBaseMixin from '@/mixins/rights/entities/base';

export default {
  mixins: [rightEntitiesBaseMixin],
  methods: {
    async createRightByPlaylistId(playlistId) {
      const checksum = USERS_RIGHTS_MASKS.delete;
      const right = {
        ...generateRight(),

        _id: playlistId,
        desc: `Rights on playlist: ${playlistId}`,
      };

      return this.createRightAndAddIntoRole(right, checksum);
    },

    async removeRightByPlaylistId(playlistId) {
      return this.removeRightAndRemoveRightFromRoleById(playlistId);
    },
  },
};
