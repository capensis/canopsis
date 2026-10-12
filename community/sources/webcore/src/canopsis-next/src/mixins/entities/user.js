import { mapActions, mapState, mapStores } from 'pinia';

import { useUserStore } from '@/stores/modules/user';

export const entitiesUserMixin = {
  computed: {
    ...mapStores(useUserStore),
    ...mapState(useUserStore, {
      users: 'items',
      usersPending: 'pending',
      usersMeta: 'meta',
    }),
  },

  methods: {
    ...mapActions(useUserStore, {
      fetchUsersList: 'fetchList',
      fetchUsersListWithPreviousParams: 'fetchListWithPreviousParams',
      createUser: 'create',
      updateUser: 'update',
      updateCurrentUser: 'updateCurrentUser',
      removeUser: 'remove',
      bulkEnableUsers: 'bulkEnable',
      bulkDisableUsers: 'bulkDisable',
      bulkRemoveUsers: 'bulkRemove',
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

export default entitiesUserMixin;
