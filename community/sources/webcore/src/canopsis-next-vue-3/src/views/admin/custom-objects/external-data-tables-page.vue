<template>
  <c-page
    :create-tooltip="t('modals.createExternalDataTable.create.title')"
    :creatable="hasCreateExternalDataTableAccess"
    @create="showCreateExternalDataTableModal"
    @refresh="refresh"
  >
    <external-data-tables-list
      :external-data-tables="externalDataTables"
      :pending="pending"
      :options="options"
      :total-items="meta.total_count"
      :updatable="hasUpdateExternalDataTableAccess"
      :removable="hasDeleteExternalDataTableAccess"
      @edit="showEditExternalDataTableModal"
      @remove="showRemoveExternalDataTableModal"
      @update:options="updateOptions"
    />
  </c-page>
</template>

<script>
import { onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useObserver } from '@/hooks/observer';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useExternalDataTable } from '@/hooks/store/modules/external-data-table';

import ExternalDataTablesList from '@/components/other/external-data-table/external-data-tables-list.vue';

export default {
  components: { ExternalDataTablesList },

  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const {
      hasCreateAccess: hasCreateExternalDataTableAccess,
      hasUpdateAccess: hasUpdateExternalDataTableAccess,
      hasDeleteAccess: hasDeleteExternalDataTableAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.externalDataTable);

    /**
     * STORE
     */
    const {
      createExternalDataTable,
      updateExternalDataTable,
      removeExternalDataTable,
      fetchExternalDataTablesListWithoutStore,
    } = useExternalDataTable();

    /**
     * QUERY
     */
    const {
      data: externalDataTables,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: ({ params } = {}) =>
        fetchExternalDataTablesListWithoutStore({
          params: {
            ...params,
            with_flags: true,
          },
        }),
    });

    /**
     * Displays a modal for creating a new external data table.
     */
    const showCreateExternalDataTableModal = () =>
      modals.show({
        name: MODALS.createExternalDataTable,
        config: {
          action: async (data) => {
            await createExternalDataTable({ data });

            return fetchList();
          },
        },
      });

    /**
     * Displays a modal for editing an existing external data table.
     *
     * @param {Object} externalDataTable - The external data table to be edited.
     */
    const showEditExternalDataTableModal = externalDataTable =>
      modals.show({
        name: MODALS.createExternalDataTable,
        config: {
          externalDataTable,
          action: async (data) => {
            await updateExternalDataTable({ id: externalDataTable._id, data });

            return fetchList();
          },
        },
      });

    /**
     * Displays a confirmation modal for removing an external data table.
     *
     * @param {Object} externalDataTable - The external data table to be removed.
     */
    const showRemoveExternalDataTableModal = externalDataTable =>
      modals.show({
        name: MODALS.confirmationPhrase,
        config: {
          title: t('modals.confirmationPhrase.deleteExternalDataTable.title'),
          text: t('modals.confirmationPhrase.deleteExternalDataTable.text'),
          phraseText: t('modals.confirmationPhrase.deleteExternalDataTable.phraseText'),
          phrase: externalDataTable.name,
          action: async () => {
            await removeExternalDataTable({ id: externalDataTable._id });

            return fetchList();
          },
        },
      });

    const { observer } = useObserver({ key: '$refresh' });

    const refresh = () => observer.notify();

    onMounted(() => {
      observer.register(fetchList);
      fetchList();
    });

    return {
      t,
      hasCreateExternalDataTableAccess,
      hasUpdateExternalDataTableAccess,
      hasDeleteExternalDataTableAccess,

      externalDataTables,
      meta,
      pending,
      options,

      refresh,
      updateOptions,
      showCreateExternalDataTableModal,
      showEditExternalDataTableModal,
      showRemoveExternalDataTableModal,
    };
  },
};
</script>
