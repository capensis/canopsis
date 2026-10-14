<template>
  <c-page
    :creatable="hasCreateAccess"
    :create-tooltip="t('modals.createUser.create.title')"
    @refresh="fetchUsersListWithPreviousParams"
    @create="showCreateUserModal"
  >
    <users />
  </c-page>
</template>

<script>
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useUser } from '@/hooks/store/modules/user';

import Users from '@/components/other/users/users.vue';

export default {
  components: {
    Users,
  },

  setup() {
    const modals = useModals();
    const { t } = useI18n();
    const { fetchUsersListWithPreviousParams, createUser } = useUser();
    const { hasCreateAccess } = useCRUDPermissions(USER_PERMISSIONS.technical.user);

    const showCreateUserModal = () => {
      modals.show({
        name: MODALS.createUser,
        config: {
          action: async (data) => {
            await createUser({ data });

            await fetchUsersListWithPreviousParams();
          },
        },
      });
    };

    return {
      hasCreateAccess,
      t,
      fetchUsersListWithPreviousParams,
      showCreateUserModal,
    };
  },
};
</script>
