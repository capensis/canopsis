<template lang="pug">
  v-container
    div
      div(v-show="selected.length")
        v-btn(@click="showRemoveSelectedRolesModal", icon)
          v-icon delete
      v-data-table(
      v-model="selected",
      :headers="headers",
      :items="roles",
      :pagination.sync="pagination",
      :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
      :total-items="rolesMeta.total",
      item-key="id"
      select-all,
      )
        template(slot="items", slot-scope="props")
          tr
            td
              v-checkbox(v-model="props.selected", primary, hide-details)
            td.text-xs-center {{ props.item.crecord_name }}
            td.text-xs-center
              v-btn(@click="showEditRoleModal(props.item)", icon)
                v-icon edit
              v-btn(@click="showRemoveRoleModal(props.item._id)", icon)
                v-icon(color="red darken-4") delete
    .fab
      v-tooltip(left)
        v-btn(slot="activator", fab, dark, color="indigo", @click.stop="showCreateRoleModal")
          v-icon add
        span Add role

</template>

<script>
import isEmpty from 'lodash/isEmpty';

import entitiesRoleMixins from '@/mixins/entities/role';
import modalMixin from '@/mixins/modal/modal';
import { MODALS } from '@/constants';

export default {
  mixins: [modalMixin, entitiesRoleMixins],
  data() {
    return {
      pagination: null,
      headers: [
        {
          text: 'Name',
          value: 'crecord_name',
        },
        {
          text: 'Actions',
          value: 'actions',
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
    showRemoveRoleModal(id) {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await this.removeRole({ id });
            await this.fetchRolesListWithPreviousParams();
          },
        },
      });
    },
    showRemoveSelectedRolesModal() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await Promise.all(this.selected.map(id => this.removeRole({ id })));

            this.selected = [];
          },
        },
      });
    },
    showEditRoleModal(role) {
      this.showModal({
        name: MODALS.createRole,
        config: {
          title: this.$t('modals.editRole.title'),
          item: role,
        },
      });
    },
    showCreateRoleModal() {
      this.showModal({
        name: MODALS.createRole,
        config: {
          title: this.$t('modals.createRole.title'),
        },
      });
    },
    fetchList() {
      const {
        rowsPerPage, page, sortBy, descending,
      } = this.pagination;

      this.fetchRolesList({
        params: {
          limit: rowsPerPage,
          start: (page - 1) * rowsPerPage,
          sort: [{ property: sortBy, direction: descending ? 'DESC' : 'ASC' }],
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
