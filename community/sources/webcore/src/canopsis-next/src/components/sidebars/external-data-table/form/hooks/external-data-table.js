import { ref } from 'vue';

import { useExternalDataTable } from '@/hooks/store/modules/external-data-table';
import { usePendingHandler } from '@/hooks/query/pending';

export const useExternalDataTableColumns = () => {
  const { fetchExternalDataTableWithoutStore } = useExternalDataTable();

  const columns = ref([]);

  const { pending, handler: fetchColumns } = usePendingHandler(async (id) => {
    const externalDataTable = await fetchExternalDataTableWithoutStore({ id });

    columns.value = externalDataTable?.columns ?? [];
  });

  return {
    pending,
    columns,
    fetchColumns,
  };
};
