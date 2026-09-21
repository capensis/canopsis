<template>
  <users-list
    :options="options"
    :pending="pending"
    :removable="hasDeleteAccess"
    :total-items="meta.total_count"
    :updatable="hasUpdateAccess"
    :users="users"
    @edit="showEditUserModal"
    @remove="showRemoveUserModal"
    @refresh="fetchList"
    @update:options="updateOptions"
  />
</template>

<script>
import { onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { getQueryForList } from '@/helpers/entities/shared/query';

import { useCRUDPermissions, useAuth } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useUser } from '@/hooks/store/modules/user';

import UsersList from '@/components/other/users/users-list.vue';

export default {
  components: {
    UsersList,
  },

  setup() {
    const modals = useModals();
    const popups = usePopups();
    const { t } = useI18n();
    const { currentUser, fetchCurrentUser } = useAuth();
    const { users, usersMeta, usersPending, fetchUsersList, updateUser, removeUser } = useUser();
    const { hasUpdateAccess, hasDeleteAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.user,
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
      onUpdate: fetchQuery => fetchUsersList({
          params: {
            ...getQueryForList(fetchQuery),
            with_flags: true,
          },
        }),
    });

    const showRemoveUserModal = (user) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await removeUser({ id: user._id });

            popups.success({ text: t('success.default') });

            await fetchList();
          },
        },
      });
    };

    const showEditUserModal = user => modals.show({
        name: MODALS.createUser,
        config: {
          title: t('modals.createUser.edit.title'),
          user,
          action: async (data) => {
            await updateUser({ data, id: user._id });

            const requests = [fetchList()];
            const isCurrentUser = user._id === currentUser.value?._id;

            if (isCurrentUser) {
              requests.push(fetchCurrentUser());
            }

            await Promise.all(requests);

            popups.success({ text: t('success.default') });
          },
        },
      });

    onMounted(fetchList);

    return {
      users,
      meta: usersMeta,
      pending: usersPending,
      options,
      hasUpdateAccess,
      hasDeleteAccess,
      fetchList,
      updateOptions,
      showEditUserModal,
      showRemoveUserModal,
    };
  },
};
</script>
