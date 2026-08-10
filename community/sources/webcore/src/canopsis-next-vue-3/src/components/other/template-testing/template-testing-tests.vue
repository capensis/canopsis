<template>
  <template-testing-tests-list
    :options="options"
    :items="items"
    :pending="pending"
    :total-items="meta.total_count"
    updatable
    removable
    @edit="showEditTemplateTestingTestModal"
    @remove="showRemoveTemplateTestingTestModal"
    @update:options="updateOptions"
  />
</template>

<script>
import { onMounted } from 'vue';

import { MODALS } from '@/constants';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useTemplateTesting } from '@/hooks/store/modules/template-testing';

import TemplateTestingTestsList from './partials/template-testing-tests-list.vue';

export default {
  components: { TemplateTestingTestsList },
  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const {
      fetchTemplateTestingTestsListWithoutStore,
      updateTemplateTestingTest,
      removeTemplateTestingTest,
    } = useTemplateTesting();

    const {
      data: items,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: fetchTemplateTestingTestsListWithoutStore,
    });

    const showEditTemplateTestingTestModal = (templateTestingTest = {}) =>
      modals.show({
        name: MODALS.createTemplateTestingTest,
        config: {
          templateTestingTest,
          title: t('modals.createTemplateTestingTest.edit.title'),
          action: async (newTemplateTestingTest) => {
            await updateTemplateTestingTest({
              id: templateTestingTest._id,
              data: newTemplateTestingTest,
            });

            return fetchList();
          },
        },
      });

    const showRemoveTemplateTestingTestModal = (templateTestingTest = {}) =>
      modals.show({
        name: MODALS.confirmationPhrase,
        config: {
          phrase: templateTestingTest.name,
          title: t('modals.confirmationPhrase.templateTestingTest.title'),
          text: t('modals.confirmationPhrase.templateTestingTest.text'),
          phraseText: t('modals.confirmationPhrase.templateTestingTest.phraseText'),
          action: async () => {
            await removeTemplateTestingTest({ id: templateTestingTest._id });

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

      updateOptions,
      fetchList,

      showEditTemplateTestingTestModal,
      showRemoveTemplateTestingTestModal,
    };
  },
};
</script>
