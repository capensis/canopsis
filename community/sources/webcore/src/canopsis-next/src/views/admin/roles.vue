<template lang="pug">
  v-container
    c-the-page-header {{ $t('common.roles') }}
    roles-list(
      :roles="roles",
      :pending="rolesPending",
      :pagination.sync="pagination",
      :total-items="rolesMeta.total_count",
      @edit="showEditRoleModal",
      @remove="showRemoveRoleModal",
      @remove-selected="showRemoveSelectedRolesModal"
    )
    c-fab-btn(
      :has-access="hasCreateAnyRoleAccess",
      @refresh="fetchList",
      @create="showCreateRoleModal"
    )
      span {{ $t('modals.createRole.create.title') }}
</template>

<script>
import { MODALS } from '@/constants';

import entitiesRoleMixin from '@/mixins/entities/role';
import { permissionsTechnicalRoleMixin } from '@/mixins/permissions/technical/role';
import { localQueryMixin } from '@/mixins/query-local/query';

import RolesList from '@/components/other/roles/roles-list.vue';

export default {
  components: {
    RolesList,
  },
  mixins: [
    localQueryMixin,
    entitiesRoleMixin,
    permissionsTechnicalRoleMixin,
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
              await this.fetchList();

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

              await this.fetchList();

              this.$popups.success({ text: this.$t('success.default') });
            } catch (err) {
              this.$popups.error({ text: this.$t('errors.default') });
            }
          },
        },
      });
    },

    showEditRoleModal(role) {
      this.$modals.show({
        name: MODALS.createRole,
        config: {
          title: this.$t('modals.createRole.edit.title'),
          role,
          action: async (data) => {
            await this.updateRole({ data, id: role._id });

            this.$popups.success({ text: this.$t('success.default') });

            await this.fetchList();
          },
        },
      });
    },

    showCreateRoleModal() {
      this.$modals.show({
        name: MODALS.createRole,
        config: {
          action: async (data) => {
            await this.createRole({ data });

            this.$popups.success({ text: this.$t('success.default') });

            await this.fetchList();
          },
        },
      });
    },

    fetchList() {
      this.fetchRolesList({ params: this.getQuery() });
    },
  },
};
</script>