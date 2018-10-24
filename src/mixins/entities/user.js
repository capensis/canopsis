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
    }),
  },
  methods: {
    ...mapActions({
      fetchUsersList: 'fetchList',
      createUser: 'create',
      removeUser: 'remove',
    }),
  },
};
