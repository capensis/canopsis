<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.roles') }}
    div
      div(v-show="hasDeleteAnyRoleAccess && selected.length")
        v-btn(@click="showRemoveSelectedRolesModal", icon)
          v-icon delete
      v-data-table(
      v-model="selected",
      :headers="headers",
      :items="roles",
      :loading="rolesPending",
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
            td {{ props.item.crecord_name }}
            td
              v-btn.ma-0(v-if="hasUpdateAnyRoleAccess", @click="showEditRoleModal(props.item._id)", icon)
                v-icon edit
              v-btn.ma-0(v-if="hasDeleteAnyRoleAccess", @click="showRemoveRoleModal(props.item._id)", icon)
                v-icon(color="error") delete
    .fab(v-if="hasCreateAnyRoleAccess")
      v-layout(column)
        v-tooltip(left)
          v-btn(slot="activator", fab, dark, icon, color="secondary", @click.stop="fetchList")
            v-icon refresh
          span {{ $t('common.refresh') }}
        v-tooltip(left)
          v-btn.secondary(slot="activator", fab, dark, @click.stop="showCreateRoleModal")
            v-icon add
          span {{ $t('modals.createRole.title') }}
</template>

<script>
import isEmpty from 'lodash/isEmpty';

import { MODALS } from '@/constants';

import popupMixin from '@/mixins/popup';
import modalMixin from '@/mixins/modal';
import entitiesRoleMixins from '@/mixins/entities/role';
import rightsTechnicalRoleMixin from '@/mixins/rights/technical/role';

export default {
  mixins: [
    popupMixin,
    modalMixin,
    entitiesRoleMixins,
    rightsTechnicalRoleMixin,
  ],
  data() {
    return {
      pagination: null,
      headers: [
        {
          text: this.$t('tables.rolesList.name'),
          value: 'crecord_name',
        },
        {
          text: this.$t('tables.rolesList.actions'),
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
            try {
              await this.removeRole({ id });
              await this.fetchRolesListWithPreviousParams();

              this.addSuccessPopup({ text: this.$t('success.default') });
            } catch (err) {
              this.addErrorPopup({ text: this.$t('errors.default') });
            }
          },
        },
      });
    },

    showRemoveSelectedRolesModal() {
      this.showModal({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await Promise.all(this.selected.map(id => this.removeRole({ id })));

              this.selected = [];

              this.addSuccessPopup({ text: this.$t('success.default') });
            } catch (err) {
              this.addErrorPopup({ text: this.$t('errors.default') });
            }
          },
        },
      });
    },

    showEditRoleModal(roleId) {
      this.showModal({
        name: MODALS.createRole,
        config: {
          title: this.$t('modals.editRole.title'),
          roleId,
        },
      });
    },

    showCreateRoleModal() {
      this.showModal({
        name: MODALS.createRole,
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
