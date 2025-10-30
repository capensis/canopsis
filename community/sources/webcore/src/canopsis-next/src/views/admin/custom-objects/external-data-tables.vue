<template>
  <c-page
    :create-tooltip="$t('modals.createExternalDataTable.create.title')"
    :creatable="hasCreateExternalDataTableAccess"
    @create="showCreateExternalDataTableModal"
    @refresh="refresh"
  >
    <external-data-tables-list
      :external-data-tables="externalDataTables"
      :pending="pending"
      :options.sync="options"
      :total-items="meta.total_count"
      :updatable="hasUpdateExternalDataTableAccess"
      :removable="hasDeleteExternalDataTableAccess"
      @edit="showEditExternalDataTableModal"
      @remove="showRemoveExternalDataTableModal"
    />
  </c-page>
</template>

<script>
import { ref, onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useObserver } from '@/hooks/observer';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useQueryOptions } from '@/hooks/query/options';
import { useCRUDPermissions } from '@/hooks/auth';
import { useExternalDataTable } from '@/hooks/store/modules/external-data-table';

import ExternalDataTablesList from '@/components/other/external-data-table/external-data-tables-list.vue';

export default {
  components: { ExternalDataTablesList },
  setup() {
    const externalDataTables = ref([]);
    const meta = ref({});

    const { t } = useI18n();
    const modals = useModals();
    const {
      hasCreateAccess: hasCreateExternalDataTableAccess,
      hasUpdateAccess: hasUpdateExternalDataTableAccess,
      hasDeleteAccess: hasDeleteExternalDataTableAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.exploitation.externalDataTable);

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
      query,
      pending,
      updateQuery,
      handler: fetchList,
    } = usePendingWithLocalQuery({
      fetchHandler: async (fetchQuery) => {
        const response = await fetchExternalDataTablesListWithoutStore({
          params: {
            limit: fetchQuery.itemsPerPage,
            page: fetchQuery.page,
            with_flags: true,
          },
        });

        externalDataTables.value = response.data;
        meta.value = response.meta;
      },
    });

    const { options } = useQueryOptions(query, updateQuery);

    /**
     * Displays a modal for creating a new external data table.
     */
    const showCreateExternalDataTableModal = () => modals.show({
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
    const showEditExternalDataTableModal = externalDataTable => modals.show({
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
    const showRemoveExternalDataTableModal = externalDataTable => modals.show({
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
      hasCreateExternalDataTableAccess,
      hasUpdateExternalDataTableAccess,
      hasDeleteExternalDataTableAccess,

      externalDataTables,
      meta,
      pending,
      options,

      refresh,
      showCreateExternalDataTableModal,
      showEditExternalDataTableModal,
      showRemoveExternalDataTableModal,
    };
  },
};
</script>
