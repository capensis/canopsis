<template>
  <c-page
    :creatable="hasCreateAccess"
    :create-tooltip="t('modals.createTheme.create.title')"
    @refresh="fetchList"
    @create="showCreateThemeModal"
  >
    <themes-list
      :themes="themes"
      :pending="themesPending"
      :options="options"
      :total-items="themesMeta.total_count"
      :updatable="hasUpdateAccess"
      :removable="hasDeleteAccess"
      :duplicable="hasCreateAccess"
      @update:options="updateOptions"
      @edit="showEditThemeModal"
      @remove="showRemoveThemeModal"
      @duplicate="showDuplicateThemeModal"
      @remove-selected="showDeleteSelectedThemesModal"
    />
  </c-page>
</template>

<script>
import { omit } from 'lodash';
import { inject, onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { pickIds } from '@/helpers/array';
import { convertQueryToRequest } from '@/helpers/query';

import { useAuth, useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useTheme } from '@/hooks/store/modules/theme';

import ThemesList from '@/components/other/theme/themes-list.vue';

export default {
  components: { ThemesList },

  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const $system = inject('$system', null);
    const { currentUser } = useAuth();
    const {
      themes,
      themesMeta,
      themesPending,
      fetchThemesList,
      createTheme,
      updateTheme,
      removeTheme,
      bulkRemoveThemes,
    } = useTheme();
    const { hasCreateAccess, hasUpdateAccess, hasDeleteAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.profile.theme,
    );

    const fetchList = query => fetchThemesList({
        params: convertQueryToRequest(query ?? options.value),
      });
    const { options, updateOptions } = useLocalQueryWithOptions({
      onUpdate: fetchList,
    });

    const showCreateThemeModal = () => {
      modals.show({
        name: MODALS.createTheme,
        config: {
          action: async (newTheme) => {
            await createTheme({ data: newTheme });

            return fetchList();
          },
        },
      });
    };

    const showEditThemeModal = (theme) => {
      modals.show({
        name: MODALS.createTheme,
        config: {
          theme,
          title: t('modals.createTheme.edit.title'),
          action: async (newTheme) => {
            await updateTheme({ id: theme._id, data: newTheme });

            if (currentUser.value.ui_theme === theme._id) {
              $system?.setTheme(newTheme);
            }

            return fetchList();
          },
        },
      });
    };

    const showRemoveThemeModal = (id) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await removeTheme({ id });

            return fetchList();
          },
        },
      });
    };

    const showDuplicateThemeModal = (theme) => {
      modals.show({
        name: MODALS.createTheme,
        config: {
          theme: omit(theme, ['_id']),
          title: t('modals.createTheme.duplicate.title'),
          action: async (newTheme) => {
            await createTheme({ data: newTheme });

            return fetchList();
          },
        },
      });
    };

    const showDeleteSelectedThemesModal = (selected) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await bulkRemoveThemes({ data: pickIds(selected) });

            return fetchList();
          },
        },
      });
    };

    onMounted(fetchList);

    return {
      t,
      themes,
      themesMeta,
      themesPending,
      options,
      hasCreateAccess,
      hasUpdateAccess,
      hasDeleteAccess,
      updateOptions,
      fetchList,
      showCreateThemeModal,
      showEditThemeModal,
      showRemoveThemeModal,
      showDuplicateThemeModal,
      showDeleteSelectedThemesModal,
    };
  },
};
</script>
