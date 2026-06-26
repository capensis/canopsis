<template>
  <div>
    <c-page-header />

    <v-tabs
      v-model="activeTab"
      bg-color="transparent"
      class="ma-4 mt-0"
      color="primary"
      align-tabs="center"
    >
      <v-tab
        v-if="hasReadAccessForTemplateData"
        :value="TEMPLATE_TESTING_TABS.data"
      >
        {{ t('templateTesting.tabs.data') }}
      </v-tab>
      <v-tab :value="TEMPLATE_TESTING_TABS.tests">
        {{ t('templateTesting.tabs.tests') }}
      </v-tab>
    </v-tabs>

    <v-card class="mx-4 mb-4 mt-0">
      <v-window v-model="activeTab">
        <v-window-item
          v-if="hasReadAccessForTemplateData"
          :value="TEMPLATE_TESTING_TABS.data"
        >
          <template-testing-data ref="templateTestingDataElement" />
        </v-window-item>

        <v-window-item :value="TEMPLATE_TESTING_TABS.tests">
          <template-testing-tests ref="templateTestingTestsElement" />
        </v-window-item>
      </v-window>
    </v-card>

    <c-fab-btn
      :has-access="creatable"
      @refresh="refresh"
      @create="showCreateTemplateTestingDataModal"
    >
      <span>{{ t('modals.createTemplateData.title') }}</span>
    </c-fab-btn>
  </div>
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

    const creatable = computed(() => (
      activeTab.value === TEMPLATE_TESTING_TABS.data && hasCreateAccessForTemplateData.value
    ));

    const refresh = () => ({
      [TEMPLATE_TESTING_TABS.data]: templateTestingDataElement.value?.fetchList,
      [TEMPLATE_TESTING_TABS.tests]: templateTestingTestsElement.value?.fetchList,
    }[activeTab.value]?.());

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
