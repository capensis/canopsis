<template>
  <users-list
    :options="options"
    :pending="pending"
    :removable="hasDeleteAccess"
    :total-items="meta.total_count"
    :updatable="hasUpdateAccess"
    :show-top-user-by-connections="showTopUserByConnections"
    :users="shownUsers"
    @edit="showEditUserModal"
    @remove="showRemoveUserModal"
    @refresh="fetchList"
    @update:options="updateListOptions"
    @update:show-top-user-by-connections="updateShowTopUserByConnections"
  />
</template>

<script>
import { computed, onMounted, ref } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { getQueryForList } from '@/helpers/entities/shared/query';

import { useCRUDPermissions, useAuth } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { usePendingHandler } from '@/hooks/query/pending';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useUser } from '@/hooks/store/modules/user';

import UsersList from '@/components/other/users/users-list.vue';

const TOP_USER_OPTIONS = {
  page: 1,
  itemsPerPage: PAGINATION_LIMIT,
  sortBy: [],
  search: '',
};

export default {
  components: {
    UsersList,
  },

  setup() {
    const modals = useModals();
    const popups = usePopups();
    const { t } = useI18n();
    const { currentUser, fetchCurrentUser } = useAuth();
    const {
      users,
      usersMeta,
      usersPending,
      fetchUsersList,
      fetchTopUserByActiveConnectsWithoutStore,
      updateUser,
      removeUser,
    } = useUser();
    const { hasUpdateAccess, hasDeleteAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.user,
    );

    const showTopUserByConnections = ref(false);
    const topUser = ref();

    /**
     * Fetches the user with the highest number of active connections.
     *
     * @returns {Promise<Object|undefined>} The top connected user, or undefined when none exists.
     */
    const fetchTopUserHandler = async () => {
      topUser.value = await fetchTopUserByActiveConnectsWithoutStore();

      return topUser.value;
    };

    const { pending: topUserPending, handler: fetchTopUser } = usePendingHandler(
      fetchTopUserHandler,
    );

    /**
     * Fetches the regular users list using the supplied table query.
     * The request is skipped while top-user mode is active.
     *
     * @param {Object} fetchQuery - Pagination, sorting, and search options for the users list.
     * @returns {Promise<Object>|undefined} The users-list response, or undefined when skipped.
     */
    const fetchUsersByQuery = (fetchQuery) => {
      if (showTopUserByConnections.value) {
        return undefined;
      }

      return fetchUsersList({
        params: {
          ...getQueryForList(fetchQuery),
          with_flags: true,
        },
      });
    };

    const {
      options: usersOptions,
      handler: fetchUsersListWithCurrentQuery,
      updateOptions: updateUsersOptions,
    } = useLocalQueryWithOptions({
      initialQuery: {
        page: 1,
        itemsPerPage: PAGINATION_LIMIT,
        sortBy: [],
        sortDesc: [],
        search: '',
      },
      onUpdate: fetchUsersByQuery,
    });

    const shownUsers = computed(() => {
      if (!showTopUserByConnections.value) {
        return users.value;
      }

      return topUser.value ? [topUser.value] : [];
    });
    const options = computed(() => (
      showTopUserByConnections.value ? TOP_USER_OPTIONS : usersOptions.value
    ));
    const pending = computed(() => (
      showTopUserByConnections.value ? topUserPending.value : usersPending.value
    ));

    /**
     * Refreshes the data source selected by the current display mode.
     */
    const fetchList = () => (
      showTopUserByConnections.value ? fetchTopUser() : fetchUsersListWithCurrentQuery()
    );

    /**
     * Switches between the regular users list and the top connected user.
     * Enabling top-user mode resets the regular search and page without querying the regular endpoint.
     *
     * @param {boolean} value - Whether top-user mode should be enabled.
     */
    const updateShowTopUserByConnections = (value) => {
      showTopUserByConnections.value = value;

      if (value) {
        if (usersOptions.value.search || usersOptions.value.page !== 1) {
          updateUsersOptions({
            ...usersOptions.value,
            page: 1,
            search: '',
          });
        }

        return fetchTopUser();
      }

      return fetchUsersListWithCurrentQuery();
    };

    /**
     * Applies table options only when the regular users list is displayed.
     *
     * @param {Object} value - Updated pagination, sorting, and search options.
     * @returns {void}
     */
    const updateListOptions = (value) => {
      if (!showTopUserByConnections.value) {
        updateUsersOptions(value);
      }
    };

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

    onMounted(fetchUsersListWithCurrentQuery);

    return {
      shownUsers,
      meta: usersMeta,
      pending,
      options,
      showTopUserByConnections,
      hasUpdateAccess,
      hasDeleteAccess,
      fetchList,
      updateListOptions,
      updateShowTopUserByConnections,
      showEditUserModal,
      showRemoveUserModal,
    };
  },
};
</script>
