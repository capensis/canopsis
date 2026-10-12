<template>
  <c-fab-expand-btn v-if="hasCreateAnyUserAccess || hasCreateAnyRoleAccess" @refresh="refresh">
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

<script setup>
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';
import { useModals } from '@/hooks/modals';
import { useRole } from '@/hooks/store/modules/role';
import { useUser } from '@/hooks/store/modules/user';

const emit = defineEmits(['refresh']);

const modals = useModals();

const { createUserWithPopup } = useUser();
const { createRoleWithPopup } = useRole();

const { hasCreateAccess: hasCreateAnyUserAccess } = useCRUDPermissions(
  USER_PERMISSIONS.technical.user,
);

const { hasCreateAccess: hasCreateAnyRoleAccess } = useCRUDPermissions(
  USER_PERMISSIONS.technical.role,
);

const refresh = () => emit('refresh');

const showCreateUserModal = () => modals.show({
    name: MODALS.createUser,
    config: {
      action: data => createUserWithPopup({ data }),
    },
  });

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
</script>
