<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.users') }}
    div
      div(v-show="hasDeleteAnyUserAccess && selected.length")
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
      item-key="_id"
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
                v-btn(v-if="hasUpdateAnyUserAccess", @click="showEditUserModal(props.item._id)", icon)
                  v-icon edit
                v-btn(v-if="hasDeleteAnyUserAccess", @click="showRemoveUserModal(props.item._id)", icon)
                  v-icon(color="red darken-4") delete
    .fab(v-if="hasCreateAnyUserAccess")
      v-layout(column)
        refresh-btn(@click="fetchList")
      v-tooltip(left)
        v-btn(slot="activator", fab, color="primary", @click.stop="showCreateUserModal")
          v-icon add
        span {{ $t('modals.createUser.title') }}
</template>

<script>
import isEmpty from 'lodash/isEmpty';

import { MODALS } from '@/constants';

import modalMixin from '@/mixins/modal';
import entitiesUserMixin from '@/mixins/entities/user';
import rightsTechnicalUserMixin from '@/mixins/rights/technical/user';

import RefreshBtn from '@/components/other/view/refresh-btn.vue';

export default {
  components: {
    RefreshBtn,
  },
  mixins: [modalMixin, entitiesUserMixin, rightsTechnicalUserMixin],
  data() {
    return {
      pagination: null,
      headers: [
        {
          text: this.$t('tables.admin.users.columns.username'),
          value: '_id',
        },
        {
          text: this.$t('tables.admin.users.columns.role'),
          value: 'role',
        },
        {
          text: this.$t('tables.admin.users.columns.enabled'),
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
            await Promise.all(this.selected.map(({ _id }) => this.removeUser({ id: _id })));
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
