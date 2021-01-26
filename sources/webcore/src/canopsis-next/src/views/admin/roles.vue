<template lang="pug">
  v-container
    c-the-page-header {{ $t('common.roles') }}
    roles-list(
      :roles="roles",
      :pending="rolesPending",
      :pagination.sync="pagination",
      :total-items="rolesMeta.total",
      @edit="showEditRoleModal",
      @remove="showRemoveRoleModal",
      @remove-selected="showRemoveSelectedRolesModal"
    )
    c-fab-btn(
      :has-access="hasCreateAnyRoleAccess",
      @refresh="fetchList",
      @create="showCreateRoleModal"
    )
      span {{ $t('modals.createRole.title') }}
</template>

<script>
import { MODALS } from '@/constants';

import { getRolesSearchByText } from '@/helpers/search/patterns-search';

import entitiesRoleMixins from '@/mixins/entities/role';
import rightsTechnicalRoleMixin from '@/mixins/rights/technical/role';
import localQueryMixin from '@/mixins/query-local/query';

import RolesList from '@/components/other/roles/roles-list.vue';

export default {
  components: {
    RolesList,
  },
  mixins: [
    localQueryMixin,
    entitiesRoleMixins,
    rightsTechnicalRoleMixin,
  ],
  mounted() {
    this.fetchList();
  },
  methods: {
    showRemoveRoleModal(id) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await this.removeRole({ id });
              await this.fetchRolesListWithPreviousParams();

              this.$popups.success({ text: this.$t('success.default') });
            } catch (err) {
              this.$popups.error({ text: this.$t('errors.default') });
            }
          },
        },
      });
    },

    showRemoveSelectedRolesModal(selected) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            try {
              await Promise.all(selected.map(({ _id }) => this.removeRole({ id: _id })));
              await this.fetchRolesListWithPreviousParams();
              this.selected = [];

              this.$popups.success({ text: this.$t('success.default') });
            } catch (err) {
              this.$popups.error({ text: this.$t('errors.default') });
            }
          },
        },
      });
    },

    showEditRoleModal({ _id: roleId }) {
      this.$modals.show({
        name: MODALS.createRole,
        config: {
          title: this.$t('modals.editRole.title'),
          roleId,
        },
      });
    },

    showCreateRoleModal() {
      this.$modals.show({
        name: MODALS.createRole,
      });
    },

    /**
     * TODO: Should be removed after backend change query
     */
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
