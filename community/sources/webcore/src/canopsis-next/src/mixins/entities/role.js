import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('role');

export const entitiesRoleMixin = {
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
      fetchRolesListWithoutStore: 'fetchListWithoutStore',
      fetchRolesList: 'fetchList',
      removeRole: 'remove',
      createRole: 'create',
      updateRole: 'update',
    }),
  },
};
