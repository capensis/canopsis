<template>
  <c-page
    :creatable="hasCreateAccess"
    :create-tooltip="$t('modals.createTemplateData.title')"
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

import { USER_PERMISSIONS, TEMPLATE_TESTING_TABS, MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useCRUDPermissions } from '@/hooks/auth';
import { useTemplateData } from '@/hooks/store/modules/template-data';

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

    const { t } = useI18n();
    const modals = useModals();

    const {
      createTemplateData,
    } = useTemplateData();

    const {
      hasCreateAccess: hasCreateAnyTemplateDataAccess,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.templateTesting);

    const activeTab = ref(TEMPLATE_TESTING_TABS.data);

    const tooltipText = computed(() => ({
      [TEMPLATE_TESTING_TABS.data]: t('modals.createTemplateData.title'),
    })[activeTab.value]);

    const hasCreateAccess = computed(() => ({
      [TEMPLATE_TESTING_TABS.data]: hasCreateAnyTemplateDataAccess.value,
    }[activeTab.value]));

    const refresh = () => ({
      [TEMPLATE_TESTING_TABS.data]: templateTestingDataElement.value?.fetchList,
      [TEMPLATE_TESTING_TABS.tests]: templateTestingTestsElement.value?.fetchList,
    }[activeTab.value]?.());

    const showCreateTemplateTestingDataModal = () => modals.show({
      name: MODALS.createTemplateTestingData,
      config: {
        title: t('modals.createTemplateData.title'),
        action: async (newTemplateTestingData) => {
          await createTemplateData({ data: newTemplateTestingData });

          return refresh();
        },
      },
    });

    return {
      tooltipText,
      hasCreateAccess,

      templateTestingDataElement,
      templateTestingTestsElement,

      TEMPLATE_TESTING_TABS,
      activeTab,

      refresh,
      showCreateTemplateTestingDataModal,
    };
  },
};
</script>
