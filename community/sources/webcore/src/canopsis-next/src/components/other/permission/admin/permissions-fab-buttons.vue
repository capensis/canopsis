<template>
  <c-fab-expand-btn
    v-if="hasCreateAnyUserAccess || hasCreateAnyRoleAccess || hasCreateAnyPermissionAccess"
    @refresh="refresh"
  >
    <c-action-fab-btn
      v-if="hasCreateAnyUserAccess"
      :tooltip="$t('modals.createUser.create.title')"
      color="indigo"
      icon="people"
      top
      @click="showCreateUserModal"
    />
    <c-action-fab-btn
      :tooltip="$t('modals.createRole.create.title')"
      color="deep-purple"
      icon="supervised_user_circle"
      top
      @click="showCreateRoleModal"
    />
  </c-fab-expand-btn>
</template>

<script>
import { MODALS } from '@/constants';

import { entitiesRoleMixin } from '@/mixins/entities/role';
import { entitiesUserMixin } from '@/mixins/entities/user';
import { permissionsTechnicalUserMixin } from '@/mixins/permissions/technical/user';
import { permissionsTechnicalRoleMixin } from '@/mixins/permissions/technical/role';
import { permissionsTechnicalPermissionMixin } from '@/mixins/permissions/technical/permission';

export default {
  mixins: [
    entitiesRoleMixin,
    entitiesUserMixin,
    permissionsTechnicalUserMixin,
    permissionsTechnicalRoleMixin,
    permissionsTechnicalPermissionMixin,
  ],
  methods: {
    showCreateUserModal() {
      this.$modals.show({
        name: MODALS.createUser,
        config: {
          action: data => this.createUserWithPopup({ data }),
        },
      });
    },

    showCreateRoleModal() {
      this.$modals.show({
        name: MODALS.createRole,
        config: {
          withTemplate: true,
          action: async (data) => {
            await this.createRole({ data });

            this.$popups.success({ text: this.$t('success.default') });

            this.refresh();
          },
        },
      });
    },

    refresh() {
      this.$emit('refresh');
    },
  },
};
</script>
