<template>
  <template-testing-data-list
    :options="options"
    :items="items"
    :pending="pending"
    :total-items="meta.total_count"
    :updatable="hasUpdateAnyTemplateDataAccess"
    :removable="hasDeleteAnyTemplateDataAccess"
    @edit="showEditTemplateTestingDataModal"
    @remove="showRemoveTemplateTestingDataModal"
    @update:options="updateOptions"
  />
</template>

<script>
import { onMounted } from 'vue';

import { MODALS, USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useCRUDPermissions } from '@/hooks/auth';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useTemplateData } from '@/hooks/store/modules/template-data';

import TemplateTestingDataList from './template-testing-data-list.vue';

export default {
  components: { TemplateTestingDataList },
  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const {
      fetchTemplateDataListWithoutStore,
      updateTemplateData,
      removeTemplateData,
    } = useTemplateData();

    const {
      hasUpdateAccess: hasUpdateAnyTemplateDataAccess,
      hasDeleteAccess: hasDeleteAnyTemplateDataAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.templateTesting);

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

    const showEditTemplateTestingDataModal = (templateTestingData = {}) => modals.show({
      name: MODALS.createTemplateTestingData,
      config: {
        templateTestingData,
        title: t('modals.createTemplateTestingData.edit.title'),
        action: async (newTemplateTestingData) => {
          await updateTemplateData({ id: templateTestingData._id, data: newTemplateTestingData });

          return fetchList();
        },
      },
    });

    const showRemoveTemplateTestingDataModal = (templateTestingData = {}) => modals.show({
      name: MODALS.confirmationPhrase,
      config: {
        phrase: templateTestingData.name,
        title: t('modals.confirmationPhrase.templateTestingData.title'),
        text: t('modals.confirmationPhrase.templateTestingData.text'),
        phraseText: t('modals.confirmationPhrase.templateTestingData.phraseText'),
        action: async () => {
          await removeTemplateData({ id: templateTestingData._id });

          return fetchList();
        },
      },
    });

    onMounted(fetchList);

    return {
      items,
      pending,
      meta,
      options,
      hasUpdateAnyTemplateDataAccess,
      hasDeleteAnyTemplateDataAccess,

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
