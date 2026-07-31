<template>
  <c-page
    :creatable="hasCreateAnyStateSettingAccess"
    :create-tooltip="t('modals.createStateSetting.create.title')"
    @refresh="fetchList"
    @create="showCreateStateSettingModal"
  >
    <state-settings-list
      :options="options"
      :state-settings="stateSettings"
      :total-items="stateSettingsMeta.total_count"
      :pending="stateSettingsPending"
      :addable="hasCreateAnyStateSettingAccess"
      :updatable="hasUpdateAnyStateSettingAccess"
      :removable="hasDeleteAnyStateSettingAccess"
      @update:options="updateOptions"
      @edit="showEditStateSettingModal"
      @duplicate="showDuplicateStateSettingModal"
      @remove="showRemoveStateSettingModal"
    />
  </c-page>
</template>

<script>
import { onMounted } from 'vue';

import { JUNIT_STATE_SETTING_ID, MODALS, USER_PERMISSIONS } from '@/constants';

import { convertQueryToRequest } from '@/helpers/query';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useStateSetting } from '@/hooks/store/modules/state-setting';

import StateSettingsList from '@/components/other/state-setting/state-settings-list.vue';

export default {
  components: { StateSettingsList },

  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();

    const {
      stateSettings,
      stateSettingsPending,
      stateSettingsMeta,
      fetchStateSettingsList,
      createStateSetting,
      updateStateSetting,
      removeStateSetting,
    } = useStateSetting();

    const {
      hasCreateAccess: hasCreateAnyStateSettingAccess,
      hasUpdateAccess: hasUpdateAnyStateSettingAccess,
      hasDeleteAccess: hasDeleteAnyStateSettingAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.stateSetting);

    const {
      options,
      updateOptions,
      handler: fetchList,
    } = useLocalQueryWithOptions({
      onUpdate: fetchQuery =>
        fetchStateSettingsList({
          params: convertQueryToRequest(fetchQuery),
        }),
    });

    const showCreateStateSettingModal = () => {
      modals.show({
        name: MODALS.createStateSetting,
        config: {
          title: t('modals.createStateSetting.create.title'),
          action: async (newStateSetting) => {
            const result = await createStateSetting({ data: newStateSetting });

            popups.success({ text: t('modals.createStateSetting.create.success') });
            await fetchList();

            return result;
          },
        },
      });
    };

    const showEditJunitStateSettingModal = (stateSetting) => {
      modals.show({
        name: MODALS.createJunitStateSetting,
        config: {
          stateSetting,
          title: t('modals.createJunitStateSetting.edit.title'),
          action: async (data) => {
            await updateStateSetting({ data, id: stateSetting._id });

            popups.success({ text: t('modals.createJunitStateSetting.edit.success') });
            await fetchList();
          },
        },
      });
    };

    const showEditStateSettingModal = (stateSetting) => {
      if (stateSetting._id === JUNIT_STATE_SETTING_ID) {
        showEditJunitStateSettingModal(stateSetting);

        return;
      }

      modals.show({
        name: MODALS.createStateSetting,
        config: {
          stateSetting,
          title: t('modals.createStateSetting.edit.title'),
          action: async (data) => {
            const result = await updateStateSetting({ data, id: stateSetting._id });

            popups.success({ text: t('modals.createStateSetting.edit.success') });
            await fetchList();

            return result;
          },
        },
      });
    };

    const showDuplicateStateSettingModal = (stateSetting) => {
      modals.show({
        name: MODALS.createStateSetting,
        config: {
          stateSetting,
          title: t('modals.createStateSetting.duplicate.title'),
          action: async (data) => {
            const result = await createStateSetting({ data });

            popups.success({ text: t('modals.createStateSetting.duplicate.success') });
            await fetchList();

            return result;
          },
        },
      });
    };

    const showRemoveStateSettingModal = (stateSetting) => {
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await removeStateSetting({ id: stateSetting._id });

            popups.success({ text: t('modals.createStateSetting.remove.success') });
            await fetchList();
          },
        },
      });
    };

    onMounted(fetchList);

    return {
      t,
      stateSettings,
      stateSettingsPending,
      stateSettingsMeta,
      options,
      hasCreateAnyStateSettingAccess,
      hasUpdateAnyStateSettingAccess,
      hasDeleteAnyStateSettingAccess,
      fetchList,
      updateOptions,
      showCreateStateSettingModal,
      showEditStateSettingModal,
      showDuplicateStateSettingModal,
      showRemoveStateSettingModal,
    };
  },
};
</script>
