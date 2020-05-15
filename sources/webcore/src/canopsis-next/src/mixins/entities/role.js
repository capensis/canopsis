import { has } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { setField, unsetField } from '@/helpers/immutable';
import { generateRoleRightByChecksum } from '@/helpers/entities';

const { mapGetters, mapActions } = createNamespacedHelpers('role');

export default {
  computed: {
    ...mapGetters({
      roles: 'items',
      getRoleById: 'getItemById',
      rolesPending: 'pending',
      rolesMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchRoleWithoutStore: 'fetchItemWithoutStore',
      fetchRolesListWithoutStore: 'fetchListWithoutStore',
      fetchRolesList: 'fetchList',
      fetchRolesListWithPreviousParams: 'fetchListWithPreviousParams',
      removeRole: 'remove',
      createRole: 'create',
    }),

    async addRightIntoRole(rightId, checksum) {
      const role = await this.fetchRoleWithoutStore({ id: this.currentUser.role });
      const newRole = setField(role, `rights.${rightId}`, generateRoleRightByChecksum(checksum));

      return this.createRole({
        data: newRole,
      });
    },

    async removeRightFromRoleById(rightId) {
      const rightPath = `rights.${rightId}`;
      const { data: roles } = await this.fetchRolesListWithoutStore({ params: { limit: 10000 } });
      const requests = roles.reduce((acc, role) => {
        if (has(role, rightPath)) {
          const request = this.createRole({
            data: unsetField(role, `rights.${rightId}`),
          });

          acc.push(request);
        }

        return acc;
      }, []);

      return Promise.all(requests);
    },
  },
};
