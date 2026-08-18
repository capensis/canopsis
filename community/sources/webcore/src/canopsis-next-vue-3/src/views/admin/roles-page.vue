<template>
  <div>
    <c-page-header />

    <v-card class="ma-4 mt-0">
      <v-card-text>
        <roles />
      </v-card-text>
    </v-card>

    <c-fab-btn :has-access="hasCreateAccess" @refresh="refresh" @create="create">
      <span>{{ t('modals.createRole.create.title') }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useRole } from '@/hooks/store/modules/role';

import Roles from '@/components/other/role/roles.vue';

export default {
  components: {
    Roles,
  },

  setup() {
    const modals = useModals();
    const popups = usePopups();
    const { t } = useI18n();
    const { createRole, fetchRolesListWithPreviousParams } = useRole();

    const { hasCreateAccess } = useCRUDPermissions(USER_PERMISSIONS.technical.role);

    const refresh = () => fetchRolesListWithPreviousParams();

    const create = () =>
      modals.show({
        name: MODALS.createRole,
        config: {
          withTemplate: true,
          action: async (data) => {
            await createRole({ data });
            popups.success({ text: t('success.default') });

            await refresh();
          },
        },
      });

    return {
      t,
      hasCreateAccess,
      refresh,
      create,
    };
  },
};
</script>
