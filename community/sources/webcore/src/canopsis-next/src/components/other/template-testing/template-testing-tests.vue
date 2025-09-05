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
import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useTemplateTest } from '@/hooks/store/modules/template-test';

import TemplateTestingTestsList from './template-testing-tests-list.vue';

export default {
  components: { TemplateTestingTestsList },
  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const {
      fetchTemplateTestListWithoutStore,
      updateTemplateTest,
      removeTemplateTest,
    } = useTemplateTest();

    const {
      data: items,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: fetchTemplateTestListWithoutStore,
    });

    const showEditTemplateTestingTestModal = (templateTestingTest = {}) => modals.show({
      name: MODALS.createTemplateTestingTest,
      config: {
        templateTestingTest,
        title: t('modals.createTemplateTestingTest.edit.title'),
        action: async (newTemplateTestingTest) => {
          await updateTemplateTest({ id: templateTestingTest._id, data: newTemplateTestingTest });

          return fetchList();
        },
      },
    });

    const showRemoveTemplateTestingTestModal = (templateTestingTest = {}) => modals.show({
      name: MODALS.confirmationPhrase,
      config: {
        phrase: templateTestingTest.name,
        title: t('modals.confirmationPhrase.templateTestingTest.title'),
        text: t('modals.confirmationPhrase.templateTestingTest.text'),
        phraseText: t('modals.confirmationPhrase.templateTestingTest.phraseText'),
        action: async () => {
          await removeTemplateTest({ id: templateTestingTest._id });

          return fetchList();
        },
      },
    });

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

      showEditTemplateTestingTestModal,
      showRemoveTemplateTestingTestModal,
    };
  },
};
</script>
