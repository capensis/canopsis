<template lang="pug">
  v-container
    div
      div(v-show="selected.length")
        v-btn(@click="showRemoveSelectedUsersModal", icon)
          v-icon delete
      v-data-table(
      v-model="selected",
      :headers="headers",
      :items="users",
      :pagination.sync="pagination",
      :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
      :total-items="usersMeta.total",
      :loading="usersPending",
      item-key="id"
      select-all,
      )
        template(slot="items", slot-scope="props")
          tr
            td
              v-checkbox(v-model="props.selected", primary, hide-details)
            td {{ props.item.id }}
            td {{ props.item.role }}
            td
              v-checkbox(:input-value="props.item.enable", primary, hide-details, disabled)
            td
              div
                v-btn(@click="showEditUserModal(props.item.id)", icon)
                  v-icon edit
                v-btn(@click="showRemoveUserModal(props.item.id)", icon)
                  v-icon delete
    .fab
      v-tooltip(left)
        v-btn(slot="activator", fab, dark, color="indigo", @click.stop="showCreateUserModal")
          v-icon add
        span Add user
</template>

<script>
import isEmpty from 'lodash/isEmpty';

import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal/modal';
import entitiesUserMixin from '@/mixins/entities/user';

export default {
  mixins: [modalMixin, entitiesUserMixin],
  data() {
    return {
      pagination: null,
      headers: [
        {
          text: 'ID',
          value: 'id',
        },
        {
          text: 'Role',
          value: 'role',
        },
        {
          text: 'Enabled',
          value: 'enable',
        },
        {
          text: '',
          sortable: false,
        },
      ],
      selected: [],
    };
  },
  watch: {
    pagination(value, oldValue) {
      if (!isEmpty(oldValue) && value !== oldValue) {
        this.fetchList();
      }
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    showRemoveUserModal(id) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeUser({ id });
            await this.fetchUsersListWithPreviousParams();
          },
        },
      });
    },

    showRemoveSelectedUsersModal() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(this.selected.map(({ id }) => this.removeUser({ id })));
            await this.fetchUsersListWithPreviousParams();

            this.selected = [];
          },
        },
      });
    },

    showEditUserModal(id) {
      this.showModal({
        name: MODALS.createUser,
        config: {
          title: this.$t('modals.editUser.title'),
          userId: id,
        },
      });
    },

    showCreateUserModal() {
      this.showModal({
        name: MODALS.createUser,
        config: {
          title: this.$t('modals.createUser.title'),
        },
      });
    },

    fetchList() {
      const {
        rowsPerPage, page, sortBy, descending,
      } = this.pagination;

      this.fetchUsersList({
        params: {
          limit: rowsPerPage,
          start: (page - 1) * rowsPerPage,
          sort: [{
            property: sortBy,
            direction: descending ? 'DESC' : 'ASC',
          }],
        },
      });
    },
  },
};
</script>

<style scoped>
  .fab {
    position: fixed;
    bottom: 0;
    right: 0;
  }
</style>
