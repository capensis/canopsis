import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('role');

export default {
  computed: {
    ...mapGetters({
      roles: 'items',
      rolesPending: 'pending',
      rolesMeta: 'meta',
    }),
  },
  methods: {
    ...mapActions({
      fetchRolesList: 'fetchList',
      fetchRolesListWithPreviousParams: 'fetchListWithPreviousParams',
      removeRole: 'remove',
      createRole: 'create',
    }),
  },
};
