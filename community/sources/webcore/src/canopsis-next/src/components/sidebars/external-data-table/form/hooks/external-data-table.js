import { ref } from 'vue';

import { useExternalDataTable } from '@/hooks/store/modules/external-data-table';
import { usePendingHandler } from '@/hooks/query/pending';

/**
 * Provides reactive state and handler for fetching external data table columns.
 *
 * @returns {Object} An object containing:
 * - `pending`: A reactive reference indicating whether the fetch operation is in progress.
 * - `columns`: A reactive reference holding the fetched columns of the external data table.
 * - `fetchColumns`: A function to initiate the fetch operation for columns by table ID.
 */
export const useExternalDataTableColumns = () => {
  const { fetchExternalDataTableWithoutStore } = useExternalDataTable();

  const columns = ref([]);

  const { pending, handler: fetchColumns } = usePendingHandler(async (id) => {
    const externalDataTable = await fetchExternalDataTableWithoutStore({ id });

    columns.value = externalDataTable?.column_configs ?? [];
  });

  return {
    pending,
    columns,
    fetchColumns,
  };
};
