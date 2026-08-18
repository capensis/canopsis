<template>
  <c-layout gap="2" column>
    <v-tabs v-model="activeTab" align-tabs="center" color="primary">
      <v-tab value="general">
        {{ t('common.general') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item value="general">
        <external-auth-token-general-form :template-vars="templateVars" />
      </v-window-item>
    </v-window>
  </c-layout>
</template>

<script>
import { onMounted, ref } from 'vue';

import { TEMPLATE_TESTING_TEST_TYPES } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useTemplateVarsList } from '@/hooks/vars/template';

import ExternalAuthTokenGeneralForm from './external-auth-token-general-form.vue';

export default {
  components: { ExternalAuthTokenGeneralForm },

  setup() {
    const { t } = useI18n();
    const activeTab = ref('general');
    const type = TEMPLATE_TESTING_TEST_TYPES.externalAuthToken;
    const { vars: templateVars, fetchList: fetchTemplateVars } = useTemplateVarsList({ type });

    onMounted(fetchTemplateVars);

    return {
      t,
      activeTab,
      templateVars,
    };
  },
};
</script>
