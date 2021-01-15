<template lang="pug">
  div
    c-the-page-header {{ $t('common.users') }}
    v-card-text
      users-list(
        :users="users",
        :total-items="usersMeta.total",
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
      span {{ $t('modals.createUser.title') }}
</template>

<script>
import { MODALS } from '@/constants';

import { prepareUserByData } from '@/helpers/entities';
import { getUsersSearchByText } from '@/helpers/search/patterns-search';

import entitiesUserMixin from '@/mixins/entities/user';
import rightsTechnicalUserMixin from '@/mixins/rights/technical/user';
import localQueryMixin from '@/mixins/query-local/query';
import authMixin from '@/mixins/auth';

import UsersList from '@/components/other/users/users-list.vue';

export default {
  components: {
    UsersList,
  },
  mixins: [entitiesUserMixin, localQueryMixin, rightsTechnicalUserMixin, authMixin],
  mounted() {
    this.fetchList();
  },
  methods: {
    showCreateUserModal() {
      this.$modals.show({
        name: MODALS.createUser,
        config: {
          action: async (data) => {
            await this.createUserWithPopup({ data: prepareUserByData(data) });

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
          title: this.$t('modals.editUser.title'),
          user,
          action: async (data) => {
            await this.createUserWithPopup({ data: prepareUserByData(data, user) });

            const requests = [this.fetchList()];

            if (user._id === this.currentUser._id) {
              requests.push(this.fetchCurrentUser());
            }

            await Promise.all(requests);
          },
        },
      });
    },

    getQuery({
      page,
      search,
      rowsPerPage,
      sortKey,
      sortDir,
    } = this.query) {
      const query = {};

      query.limit = rowsPerPage;
      query.start = (page - 1) * rowsPerPage;

      if (sortKey) {
        query.sort = [{
          property: sortKey,
          direction: sortDir,
        }];
      }

      if (search) {
        query.filter = { $and: [getUsersSearchByText(search)] };
      }

      return query;
    },

    fetchList() {
      return this.fetchUsersList({ params: this.getQuery() });
    },
  },
};
</script>
