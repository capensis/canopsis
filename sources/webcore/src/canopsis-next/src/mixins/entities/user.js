import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('user');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      users: 'items',
      usersPending: 'pending',
      usersMeta: 'meta',
      getUserById: 'getItemById',
    }),
  },
  methods: {
    ...mapActions({
      fetchUsersList: 'fetchList',
      fetchUsersListWithPreviousParams: 'fetchListWithPreviousParams',
      createUser: 'create',
      removeUser: 'remove',
    }),
  },
};
