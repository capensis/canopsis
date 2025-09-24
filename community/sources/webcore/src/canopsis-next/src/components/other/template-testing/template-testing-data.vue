<template>
  <template-testing-data-list
    :options="options"
    :items="items"
    :pending="pending"
    :total-items="meta.total_count"
    updatable
    removable
    @edit="showEditTemplateTestingDataModal"
    @remove="showRemoveTemplateTestingDataModal"
    @update:options="updateOptions"
  />
</template>

<script>
import { onMounted } from 'vue';

import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useTemplateData } from '@/hooks/store/modules/template-data';

import { useTemplateDataModals } from './hooks/template-testing-data';
import TemplateTestingDataList from './partials/template-testing-data-list.vue';

export default {
  components: { TemplateTestingDataList },
  setup() {
    const { fetchTemplateDataListWithoutStore } = useTemplateData();

    const {
      data: items,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: fetchTemplateDataListWithoutStore,
    });

    const {
      showEditTemplateTestingDataModal,
      showRemoveTemplateTestingDataModal,
    } = useTemplateDataModals(fetchList);

    onMounted(fetchList);

    return {
      items,
      pending,
      meta,
      options,

      updateOptions,

      /**
       * We need to return it for parent component to be able to fetch the list
       */
      fetchList,

      showEditTemplateTestingDataModal,
      showRemoveTemplateTestingDataModal,
    };
  },
};
</script>
