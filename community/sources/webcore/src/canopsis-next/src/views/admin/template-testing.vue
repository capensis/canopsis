<template>
  <c-page
    :create-tooltip="$t('modals.createTemplateData.title')"
    :creatable="creatable"
    @refresh="refresh"
    @create="showCreateTemplateTestingDataModal"
  >
    <v-tabs
      v-model="activeTab"
      slider-color="primary"
      fixed-tabs
    >
      <v-tab :href="`#${TEMPLATE_TESTING_TABS.data}`">
        {{ $t('templateTesting.tabs.data') }}
      </v-tab>
      <v-tab-item :value="TEMPLATE_TESTING_TABS.data">
        <v-card-text>
          <template-testing-data ref="templateTestingDataElement" />
        </v-card-text>
      </v-tab-item>
      <v-tab :href="`#${TEMPLATE_TESTING_TABS.tests}`">
        {{ $t('templateTesting.tabs.tests') }}
      </v-tab>
      <v-tab-item :value="TEMPLATE_TESTING_TABS.tests">
        <v-card-text>
          <template-testing-tests ref="templateTestingTestsElement" />
        </v-card-text>
      </v-tab-item>
    </v-tabs>
  </c-page>
</template>

<script>
import { ref, computed } from 'vue';

import { TEMPLATE_TESTING_TABS } from '@/constants';

import { useTemplateDataModals } from '@/components/other/template-testing/hooks/template-testing-data';

import TemplateTestingData from '@/components/other/template-testing/template-testing-data.vue';
import TemplateTestingTests from '@/components/other/template-testing/template-testing-tests.vue';

export default {
  components: {
    TemplateTestingData,
    TemplateTestingTests,
  },
  setup() {
    const templateTestingDataElement = ref(null);
    const templateTestingTestsElement = ref(null);

    const activeTab = ref(TEMPLATE_TESTING_TABS.data);

    const creatable = computed(() => ({ [TEMPLATE_TESTING_TABS.data]: true }[activeTab.value]));

    /**
     * Refreshes the currently active tab's data by calling the appropriate fetchList method
     */
    const refresh = () => ({
      [TEMPLATE_TESTING_TABS.data]: templateTestingDataElement.value?.fetchList,
      [TEMPLATE_TESTING_TABS.tests]: templateTestingTestsElement.value?.fetchList,
    }[activeTab.value]?.());

    const { showCreateTemplateTestingDataModal } = useTemplateDataModals({
      refresh,
    });

    return {
      templateTestingDataElement,
      templateTestingTestsElement,

      TEMPLATE_TESTING_TABS,
      activeTab,

      creatable,

      refresh,
      showCreateTemplateTestingDataModal,
    };
  },
};
</script>
