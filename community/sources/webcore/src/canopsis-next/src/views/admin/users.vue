<template lang="pug">
  div
    c-the-page-header {{ $t('common.users') }}
    v-card-text
      users-list(
        :users="users",
        :total-items="usersMeta.total_count",
        :pagination.sync="pagination",
        :pending="usersPending",
        @edit="showEditUserModal",
        @remove="showRemoveUserModal",
        @remove-selected="showRemoveSelectedUsersModal"
      )
    c-fab-btn(
      :has-access="hasCreateAnyUserAccess",
      @refresh="fetchList",
      @create="showCreateUserModal"
    )
      span {{ $t('modals.createUser.create.title') }}
</template>

<script>
import { MODALS } from '@/constants';

import entitiesUserMixin from '@/mixins/entities/user';
import { permissionsTechnicalUserMixin } from '@/mixins/permissions/technical/user';
import { localQueryMixin } from '@/mixins/query-local/query';
import { authMixin } from '@/mixins/auth';

import UsersList from '@/components/other/users/users-list.vue';

export default {
  components: {
    UsersList,
  },
  mixins: [entitiesUserMixin, localQueryMixin, permissionsTechnicalUserMixin, authMixin],
  mounted() {
    this.fetchList();
  },
  methods: {
    showCreateUserModal() {
      this.$modals.show({
        name: MODALS.createUser,
        config: {
          action: async (data) => {
            await this.createUserWithPopup({ data });

            await this.fetchList();
          },
        },
      });
    },

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

            if (user._id === this.currentUser._id) {
              requests.push(this.fetchCurrentUser());
            }

            await Promise.all(requests);
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