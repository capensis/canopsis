<template>
  <users-list
    :users="users"
    :total-items="usersMeta.total_count"
    :options.sync="options"
    :pending="usersPending"
    :removable="hasDeleteAnyUserAccess"
    :updatable="hasUpdateAnyUserAccess"
    @edit="showEditUserModal"
    @remove="showRemoveUserModal"
    @remove-selected="showRemoveSelectedUsersModal"
  />
</template>

<script>
import { MODALS } from '@/constants';

import { entitiesUserMixin } from '@/mixins/entities/user';
import { permissionsTechnicalUserMixin } from '@/mixins/permissions/technical/user';
import { localQueryMixin } from '@/mixins/query-local/query';

import UsersList from '@/components/other/users/users-list.vue';

export default {
  inject: ['$system'],
  components: {
    UsersList,
  },
  mixins: [entitiesUserMixin, localQueryMixin, permissionsTechnicalUserMixin],
  mounted() {
    this.fetchList();
  },
  methods: {
    showRemoveUserModal(user) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeUserWithPopup({ id: user._id });

            await this.fetchList();
          },
        },
      });
    },

    showRemoveSelectedUsersModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(selected.map(({ _id }) => this.removeUser({ id: _id })));

            this.$popups.success({ text: this.$t('success.default') });

            await this.fetchList();
          },
        },
      });
    },

    showEditUserModal(user) {
      this.$modals.show({
        name: MODALS.createUser,
        config: {
          title: this.$t('modals.createUser.edit.title'),
          user,
          action: async (data) => {
            await this.updateUserWithPopup({ data, id: user._id });

            const requests = [this.fetchList()];

            const isCurrentUser = user._id === this.currentUser._id;

            if (isCurrentUser) {
              requests.push(this.fetchCurrentUser());
            }

            await Promise.all(requests);

            if (isCurrentUser) {
              this.$system.setTheme(this.currentUser.ui_theme);
            }
          },
        },
      });
    },

    fetchList() {
      return this.fetchUsersList({ params: this.getQuery() });
    },
  },
};
</script>
