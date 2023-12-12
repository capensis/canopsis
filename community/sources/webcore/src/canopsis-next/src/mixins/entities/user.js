import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('user');

/**
 * @mixin
 */
export const entitiesUserMixin = {
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
      fetchUsersListWithPreviousParams: 'fetchListWithPreviousParams',
      createUser: 'create',
      updateUser: 'update',
      updateCurrentUser: 'updateCurrentUser',
      removeUser: 'remove',
    }),

    async createUserWithPopup({ data }) {
      await this.createUser({ data });

      this.$popups.success({ text: this.$t('success.default') });
    },

    async updateUserWithPopup({ data, id }) {
      await this.updateUser({ data, id });

      this.$popups.success({ text: this.$t('success.default') });
    },

    async removeUserWithPopup({ id }) {
      try {
        await this.removeUser({ id });
        this.$popups.success({ text: this.$t('success.default') });
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('errors.default') });
      }
    },
  },
};
