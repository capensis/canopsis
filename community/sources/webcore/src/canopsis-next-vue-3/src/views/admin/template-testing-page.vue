<template>
  <c-page
    :creatable="creatable"
    :create-tooltip="t('modals.createTemplateData.title')"
    @refresh="refresh"
    @create="showCreateTemplateTestingDataModal"
  >
    <v-tabs
v-model="activeTab"
bg-color="transparent"
color="primary"
align-tabs="center">
      <v-tab v-if="hasReadAccessForTemplateData" :value="TEMPLATE_TESTING_TABS.data">
        {{ t('templateTesting.tabs.data') }}
      </v-tab>
      <v-tab :value="TEMPLATE_TESTING_TABS.tests">
        {{ t('templateTesting.tabs.tests') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item v-if="hasReadAccessForTemplateData" :value="TEMPLATE_TESTING_TABS.data">
        <v-card-text>
          <template-testing-data ref="templateTestingDataElement" />
        </v-card-text>
      </v-window-item>

      <v-window-item :value="TEMPLATE_TESTING_TABS.tests">
        <v-card-text>
          <template-testing-tests ref="templateTestingTestsElement" />
        </v-card-text>
      </v-window-item>
    </v-window>
  </c-page>
</template>

<script>
import { computed, ref } from 'vue';

import { TEMPLATE_TESTING_TABS, USER_PERMISSIONS } from '@/constants';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';

import { useTemplateDataModals } from '@/components/other/template-testing/hooks/template-testing-data';
import TemplateTestingData from '@/components/other/template-testing/template-testing-data.vue';
import TemplateTestingTests from '@/components/other/template-testing/template-testing-tests.vue';

export default {
  components: {
    TemplateTestingData,
    TemplateTestingTests,
  },

  setup() {
    const { t } = useI18n();
    const templateTestingDataElement = ref(null);
    const templateTestingTestsElement = ref(null);

    const {
      hasCreateAccess: hasCreateAccessForTemplateData,
      hasReadAccess: hasReadAccessForTemplateData,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.templateData);

    const activeTab = ref(TEMPLATE_TESTING_TABS.data);

    const creatable = computed(
      () => activeTab.value === TEMPLATE_TESTING_TABS.data && hasCreateAccessForTemplateData.value,
    );

    const refresh = () =>
      ({
        [TEMPLATE_TESTING_TABS.data]: templateTestingDataElement.value?.fetchList,
        [TEMPLATE_TESTING_TABS.tests]: templateTestingTestsElement.value?.fetchList,
      })[activeTab.value]?.();

    const { showCreateTemplateTestingDataModal } = useTemplateDataModals({ refresh });

    return {
      t,
      templateTestingDataElement,
      templateTestingTestsElement,

      TEMPLATE_TESTING_TABS,
      activeTab,

      hasReadAccessForTemplateData,

      creatable,

      refresh,
      showCreateTemplateTestingDataModal,
    };
  },
};
</script>
