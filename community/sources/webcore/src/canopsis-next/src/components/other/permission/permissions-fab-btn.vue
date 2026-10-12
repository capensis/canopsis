<template>
  <c-fab-expand-btn
    v-if="hasCreateAnyUserAccess || hasCreateAnyRoleAccess"
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
      v-if="hasCreateAnyRoleAccess"
      :tooltip="$t('modals.createRole.create.title')"
      color="deep-purple"
      icon="supervised_user_circle"
      top
      @click="showCreateRoleModal"
    />
  </c-fab-expand-btn>
</template>

<script>
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { entitiesRoleMixin } from '@/mixins/entities/role';
import { entitiesUserMixin } from '@/mixins/entities/user';
import { permissionsTechnicalUserMixin } from '@/mixins/permissions/technical/user';
import { permissionsTechnicalRoleMixin } from '@/mixins/permissions/technical/role';
import { permissionsTechnicalPermissionMixin } from '@/mixins/permissions/technical/permission';
import { useUser } from '@/hooks/store/modules/user';
import { useRole } from '@/hooks/store/modules/role';
import { useModals } from '@/hooks/modals';
import { useCRUDPermissions } from '@/hooks/auth';

export default {
  mixins: [
    entitiesRoleMixin,
    entitiesUserMixin,
    permissionsTechnicalUserMixin,
    permissionsTechnicalRoleMixin,
    permissionsTechnicalPermissionMixin,
  ],
  setup(props, { emit }) {
    const modals = useModals();

    const { createUserWithPopup } = useUser();
    const { createRoleWithPopup } = useRole();

    const { hasCreateAccess: hasCreateAnyUserAccess } = useCRUDPermissions(USER_PERMISSIONS.technical.user);
    const { hasCreateAccess: hasCreateAnyRoleAccess } = useCRUDPermissions(USER_PERMISSIONS.technical.role);

    /**
     * Emits a 'refresh' event to the parent component.
     *
     * @returns {void}
     */
    const refresh = () => emit('refresh');

    /**
     * Shows the create user modal dialog.
     *
     * @returns {void}
     * @description Opens a modal for creating a new user with a configured action handler.
     */
    const showCreateUserModal = () => modals.show({
      name: MODALS.createUser,
      config: {
        action: data => createUserWithPopup({ data }),
      },
    });

    /**
     * Shows the create role modal dialog.
     *
     * @returns {void}
     * @description Opens a modal for creating a new role with template support.
     * After successful role creation, triggers a refresh.
     */
    const showCreateRoleModal = () => modals.show({
      name: MODALS.createRole,
      config: {
        withTemplate: true,
        action: async (data) => {
          await createRoleWithPopup({ data });

          refresh();
        },
      },
    });

    return {
      hasCreateAnyUserAccess,
      hasCreateAnyRoleAccess,

      refresh,
      showCreateUserModal,
      showCreateRoleModal,
    };
  },
};
</script>
