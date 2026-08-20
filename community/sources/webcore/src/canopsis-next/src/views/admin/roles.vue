<template>
  <div>
    <c-page-header />
    <v-card class="ma-4 mt-0">
      <roles-list
        :roles="roles"
        :pending="rolesPending"
        :options.sync="options"
        :total-items="rolesMeta.total_count"
        :removable="hasDeleteAnyRoleAccess"
        :duplicable="hasCreateAnyRoleAccess"
        :updatable="hasUpdateAnyRoleAccess"
        @refresh="fetchList"
        @edit="showEditRoleModal"
        @remove="showRemoveRoleModal"
        @duplicate="showDuplicateRoleModal"
        @remove-selected="showRemoveSelectedRolesModal"
      />
    </v-card>
    <c-fab-btn
      :has-access="hasCreateAnyRoleAccess"
      @refresh="fetchList"
      @create="showCreateRoleModal"
    >
      <span>{{ $t('modals.createRole.create.title') }}</span>
    </c-fab-btn>
  </div>
</template>

<script>
import { omit } from 'lodash';
import { inject, onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { convertQueryToRequest } from '@/helpers/query';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useAuth, useCRUDPermissions } from '@/hooks/auth';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useRole } from '@/hooks/store/modules/role';

import RolesList from '@/components/other/role/roles-list.vue';

export default {
  components: {
    RolesList,
  },
  setup() {
    const {
      roles,
      rolesPending,
      rolesMeta,
      fetchRolesList,
      removeRole,
      createRole,
      updateRole,
    } = useRole();

    const {
      hasCreateAccess: hasCreateAnyRoleAccess,
      hasUpdateAccess: hasUpdateAnyRoleAccess,
      hasDeleteAccess: hasDeleteAnyRoleAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.role);

    const { currentUser, fetchCurrentUser } = useAuth();
    const modals = useModals();
    const popups = usePopups();
    const { t } = useI18n();
    const system = inject('$system');

    const {
      options,
      handler: fetchList,
    } = useLocalQueryWithOptions({
      onUpdate: fetchQuery => fetchRolesList({
        params: {
          ...convertQueryToRequest(fetchQuery),
          with_flags: true,
        },
      }),
    });

    /**
     * Shows the confirmation modal for deleting a role.
     * After successful deletion, refreshes the roles list and shows a success popup.
     *
     * @param {string} id - The unique identifier of the role to delete.
     */
    const showRemoveRoleModal = id => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          try {
            await removeRole({ id });
            await fetchList();

            popups.success({ text: t('success.default') });
          } catch (err) {
            console.error(err);

            popups.error({ text: err.error ?? t('errors.default') });
          }
        },
      },
    });

    /**
     * Shows the confirmation modal for deleting multiple selected roles.
     * After successful deletion, refreshes the roles list and shows a success popup.
     *
     * @param {Object[]} selected - The list of roles to delete.
     * @param {string} selected[]._id - The unique identifier of each role.
     */
    const showRemoveSelectedRolesModal = selected => modals.show({
      name: MODALS.confirmation,
      config: {
        action: async () => {
          try {
            await Promise.all(selected.map(({ _id }) => removeRole({ id: _id })));

            await fetchList();

            popups.success({ text: t('success.default') });
          } catch (err) {
            console.error(err);

            popups.error({ text: err.error ?? t('errors.default') });
          }
        },
      },
    });

    /**
     * Shows the modal for editing an existing role.
     * Refreshes the current user and UI theme when the edited role belongs to the current user.
     *
     * @param {Object} role - The role object to edit.
     * @param {string} role._id - The unique identifier of the role.
     */
    const showEditRoleModal = role => modals.show({
      name: MODALS.createRole,
      config: {
        title: t('modals.createRole.edit.title'),
        role,
        action: async (data) => {
          await updateRole({ data, id: role._id });

          const requests = [fetchList()];

          if (currentUser.value.roles.find(currentRole => currentRole._id === role._id)) {
            requests.push(fetchCurrentUser());
          }

          popups.success({ text: t('success.default') });

          await Promise.all(requests);

          if (requests.length > 1) {
            system.setTheme(currentUser.value.ui_theme_colors);
          }
        },
      },
    });

    /**
     * Shows the modal for duplicating an existing role.
     * Creates a new role based on the provided role and refreshes the roles list.
     *
     * @param {Object} role - The role object to duplicate.
     */
    const showDuplicateRoleModal = role => modals.show({
      name: MODALS.createRole,
      config: {
        role: omit(role, ['_id']),
        title: t('modals.createRole.duplicate.title'),
        action: async (data) => {
          await createRole({ data });

          popups.success({ text: t('success.default') });

          return fetchList();
        },
      },
    });

    /**
     * Shows the modal for creating a new role.
     * Supports role templates and refreshes the roles list after successful creation.
     */
    const showCreateRoleModal = () => modals.show({
      name: MODALS.createRole,
      config: {
        withTemplate: true,
        action: async (data) => {
          await createRole({ data });

          popups.success({ text: t('success.default') });

          return fetchList();
        },
      },
    });

    onMounted(fetchList);

    return {
      roles,
      rolesPending,
      rolesMeta,
      options,
      hasCreateAnyRoleAccess,
      hasUpdateAnyRoleAccess,
      hasDeleteAnyRoleAccess,
      fetchList,
      showRemoveRoleModal,
      showRemoveSelectedRolesModal,
      showEditRoleModal,
      showDuplicateRoleModal,
      showCreateRoleModal,
    };
  },
};
</script>
