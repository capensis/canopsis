<template>
  <roles-list
    :options="options"
    :pending="pending"
    :removable="hasDeleteAccess"
    :duplicable="hasCreateAccess"
    :total-items="meta.total_count"
    :updatable="hasUpdateAccess"
    :roles="roles"
    @edit="showEditRoleModal"
    @duplicate="showDuplicateRoleModal"
    @remove="showRemoveRoleModal"
    @refresh="fetchList"
    @update:options="updateOptions"
  />
</template>

<script>
import { omit } from 'lodash';
import { onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { getQueryForList } from '@/helpers/entities/shared/query';

import { useCRUDPermissions, useAuth } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useRole } from '@/hooks/store/modules/role';

import RolesList from './roles-list.vue';

export default {
  components: {
    RolesList,
  },

  setup() {
    const modals = useModals();
    const popups = usePopups();
    const { t } = useI18n();
    const { currentUser, fetchCurrentUser } = useAuth();
    const { roles, rolesMeta, rolesPending, fetchRolesList, createRole, updateRole, removeRole } =
      useRole();

    const { hasCreateAccess, hasUpdateAccess, hasDeleteAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.role,
    );

    const {
      options,
      handler: fetchList,
      updateOptions,
    } = useLocalQueryWithOptions({
      initialQuery: {
        page: 1,
        itemsPerPage: PAGINATION_LIMIT,
        sortBy: [],
        sortDesc: [],
        search: '',
      },
      onUpdate: fetchQuery => fetchRolesList({
          params: {
            ...getQueryForList(fetchQuery),
            with_flags: true,
          },
        }),
    });

    const showRemoveRoleModal = role => modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await removeRole({ id: role._id });

            popups.success({ text: t('success.default') });

            await fetchList();
          },
        },
      });

    const showEditRoleModal = role => modals.show({
        name: MODALS.createRole,
        config: {
          title: t('modals.createRole.edit.title'),
          role,
          action: async (data) => {
            await updateRole({ data, id: role._id });

            const requests = [fetchList()];
            const currentUserRoles = currentUser.value?.roles ?? [];

            if (currentUserRoles.some(currentRole => currentRole._id === role._id)) {
              requests.push(fetchCurrentUser());
            }

            await Promise.all(requests);

            popups.success({ text: t('success.default') });
          },
        },
      });

    const showDuplicateRoleModal = role => modals.show({
        name: MODALS.createRole,
        config: {
          title: t('modals.createRole.duplicate.title'),
          role: omit(role, ['_id']),
          action: async (data) => {
            await createRole({ data });

            popups.success({ text: t('success.default') });

            await fetchList();
          },
        },
      });

    onMounted(fetchList);

    return {
      roles,
      meta: rolesMeta,
      pending: rolesPending,
      options,
      hasCreateAccess,
      hasUpdateAccess,
      hasDeleteAccess,
      updateOptions,
      showEditRoleModal,
      showDuplicateRoleModal,
      showRemoveRoleModal,
    };
  },
};
</script>
