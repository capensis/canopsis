<template>
  <template-testing-data-list
    :options="options"
    :items="items"
    :pending="pending"
    :total-items="meta.total_count"
    :updatable="hasUpdateAccessForTemplateData"
    :removable="hasDeleteAccessForTemplateData"
    @edit="showEditTemplateTestingDataModal"
    @remove="showRemoveTemplateTestingDataModal"
    @update:options="updateOptions"
  />
</template>

<script>
import { onMounted } from 'vue';

import { USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useTemplateTesting } from '@/hooks/store/modules/template-testing';

import { useTemplateDataModals } from './hooks/template-testing-data';
import TemplateTestingDataList from './partials/template-testing-data-list.vue';

export default {
  components: { TemplateTestingDataList },
  setup() {
    const {
      hasUpdateAccess: hasUpdateAccessForTemplateData,
      hasDeleteAccess: hasDeleteAccessForTemplateData,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.templateData);

    const { fetchTemplateTestingDataListWithoutStore } = useTemplateTesting();

    const {
      data: items,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: fetchTemplateTestingDataListWithoutStore,
    });

    const {
      showEditTemplateTestingDataModal,
      showRemoveTemplateTestingDataModal,
    } = useTemplateDataModals({ refresh: fetchList });

    onMounted(fetchList);

    return {
      hasUpdateAccessForTemplateData,
      hasDeleteAccessForTemplateData,

      items,
      pending,
      meta,
      options,

      updateOptions,
      fetchList,

      showEditTemplateTestingDataModal,
      showRemoveTemplateTestingDataModal,
    };
  },
};
</script>
