<template lang="pug">
  v-container
    h2.text-xs-center.my-3.display-1.font-weight-medium {{ $t('common.roles') }}
    div.white
      v-layout(row, wrap)
        v-flex(xs4)
          search-field(
          v-model="searchingText",
          @submit="applySearchFilter",
          @clear="applySearchFilter",
          )
        v-flex(v-show="hasDeleteAnyRoleAccess && selected.length", xs4)
          v-btn(@click="showRemoveSelectedRolesModal", data-test="massDeleteButton", icon)
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
          tr(:data-test="`role-${props.item._id}`")
            td
              v-checkbox(v-model="props.selected", data-test="optionCheckbox", primary, hide-details)
            td {{ props.item._id }}
            td
              v-btn.ma-0(
              v-if="hasUpdateAnyRoleAccess",
              data-test="editButton",
              icon,
              @click="showEditRoleModal(props.item._id)"
              )
                v-icon edit
              v-btn.ma-0(
              v-if="hasDeleteAnyRoleAccess",
              data-test="deleteButton",
              icon,
              @click="showRemoveRoleModal(props.item._id)"
              )
                v-icon(color="error") delete
    .fab(v-if="hasCreateAnyRoleAccess")
      v-layout(column)
        refresh-btn(@click="fetchList")
        v-tooltip(left)
          v-btn(
          slot="activator",
          color="primary",
          data-test="addButton",
          fab,
          @click.stop="showCreateRoleModal"
          )
            v-icon add
          span {{ $t('modals.createRole.title') }}
</template>

<script>
import { MODALS } from '@/constants';

import { getRolesSearchByText } from '@/helpers/entities-search';

import popupMixin from '@/mixins/popup';
import modalMixin from '@/mixins/modal';
import viewQuery from '@/mixins/view/query';
import entitiesRoleMixins from '@/mixins/entities/role';
import rightsTechnicalRoleMixin from '@/mixins/rights/technical/role';

import RefreshBtn from '@/components/other/view/buttons/refresh-btn.vue';
import SearchField from '@/components/forms/fields/search-field.vue';

export default {
  components: {
    RefreshBtn,
    SearchField,
  },
  mixins: [
    popupMixin,
    modalMixin,
    viewQuery,
    entitiesRoleMixins,
    rightsTechnicalRoleMixin,
  ],
  data() {
    return {
      searchingText: '',
      selected: [],
    };
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('tables.rolesList.name'),
          value: '_id',
        },
        {
          text: this.$t('tables.rolesList.actions'),
          value: 'actions',
        },
      ];
    },
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
              await Promise.all(this.selected.map(({ _id }) => this.removeRole({ id: _id })));
              await this.fetchRolesListWithPreviousParams();
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


    applySearchFilter() {
      this.query = {
        ...this.query,

        search: this.searchingText,
      };
    },

    getQuery() {
      const { search } = this.query;
      const query = this.getBaseQuery();

      if (search) {
        query.filter = { $and: [getRolesSearchByText(search)] };
      }

      return query;
    },

    fetchList() {
      this.fetchRolesList({ params: this.getQuery() });
    },
  },
};
</script>
