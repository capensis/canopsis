<template lang="pug">
  c-fab-expand-btn(
    v-if="hasCreateAnyUserAccess || hasCreateAnyRoleAccess || hasCreateAnyPermissionAccess",
    @refresh="refresh"
  )
    v-tooltip(v-if="hasCreateAnyUserAccess", top)
      v-btn(slot="activator", fab, dark, small, color="indigo", @click.stop="showCreateUserModal")
        v-icon people
      span {{ $t('modals.createUser.create.title') }}
    v-tooltip(v-if="hasCreateAnyRoleAccess", top)
      v-btn(slot="activator", fab, dark, small, color="deep-purple", @click.stop="showCreateRoleModal")
        v-icon supervised_user_circle
      span {{ $t('modals.createRole.create.title') }}
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
        action: async (data) => {
          await this.createRole({ data });

          this.$popups.success({ text: this.$t('success.default') });
        },
      });
    },

    refresh() {
      this.$emit('refresh');
    },
  },
};
</script>
