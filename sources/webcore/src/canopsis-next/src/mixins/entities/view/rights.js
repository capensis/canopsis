import { USERS_RIGHTS_MASKS, USERS_RIGHTS_TYPES } from '@/constants';
import { generateRight, generateRoleRightByChecksum } from '@/helpers/entities';
import { omit } from 'lodash';

import authMixin from '@/mixins/auth';
import entitiesRoleMixin from '@/mixins/entities/role';
import entitiesRightMixin from '@/mixins/entities/right';

export default {
  mixins: [authMixin, entitiesRoleMixin, entitiesRightMixin],
  methods: {
    async createRightByViewId(viewId) {
      try {
        const checksum = USERS_RIGHTS_MASKS.read + USERS_RIGHTS_MASKS.update + USERS_RIGHTS_MASKS.delete;
        const role = await this.fetchRoleWithoutStore({ id: this.currentUser.role });

        const right = {
          ...generateRight(),

          _id: viewId,
          type: USERS_RIGHTS_TYPES.rw,
          desc: `Rights on view: ${viewId}`,
        };

        await this.createRight({ data: right });
        await this.createRole({
          data: {
            ...role,
            rights: {
              ...role.rights,

              [right._id]: generateRoleRightByChecksum(checksum),
            },
          },
        });

        return this.fetchCurrentUser();
      } catch (err) {
        this.$popups.error({ text: this.$t('modals.view.errors.rightCreating') });

        return Promise.resolve();
      }
    },

    async removeRightByViewId(viewId) {
      try {
        const { data: roles } = await this.fetchRolesListWithoutStore({ params: { limit: 10000 } });

        return Promise.all([
          this.removeRight({ id: viewId }),

          ...roles.map(role => this.createRole({
            data: {
              ...role,
              rights: omit(role.rights, [viewId]),
            },
          })),
        ]);
      } catch (err) {
        this.$popups.error({ text: this.$t('modals.view.errors.rightRemoving') });

        return Promise.resolve();
      }
    },
  },
};
