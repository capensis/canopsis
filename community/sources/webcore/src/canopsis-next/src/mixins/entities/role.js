import { createNamespacedHelpers } from 'vuex';

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
      removeRole: 'remove',
      createRole: 'create',
      updateRole: 'update',
    }),
  },
};
