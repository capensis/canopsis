<template>
  <c-layout column>
    <v-tabs v-model="activeTab" align-tabs="center" color="primary">
      <v-tab :class="{ 'text-error': hasGeneralError }" value="general">
        {{ t('common.general') }}
      </v-tab>
      <v-tab :class="{ 'text-error': hasNotificationsError }" value="notifications">
        {{ t('common.notification', 2) }}
      </v-tab>
      <v-tab :class="{ 'text-error': hasAdvancedError }" value="advanced">
        {{ t('userInterface.advanced') }}
      </v-tab>
      <v-tab :class="{ 'text-error': hasLoginPageError }" value="login">
        {{ t('userInterface.loginPage') }}
      </v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item value="general" class="pt-4" eager>
        <user-interface-general-form ref="generalElement" :disabled="disabled" />
      </v-window-item>

      <v-window-item value="notifications" class="pt-4" eager>
        <user-interface-notifications-form ref="notificationsElement" :disabled="disabled" />
      </v-window-item>

      <v-window-item value="advanced" class="pt-4" eager>
        <user-interface-advanced-form ref="advancedElement" :disabled="disabled" />
      </v-window-item>

      <v-window-item value="login" class="pt-4" eager>
        <user-interface-login-form ref="loginPageElement" :disabled="disabled" />
      </v-window-item>
    </v-window>
  </c-layout>
</template>

<script>
import { ref } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';

import UserInterfaceAdvancedForm from './user-interface-advanced-form.vue';
import UserInterfaceGeneralForm from './user-interface-general-form.vue';
import UserInterfaceLoginForm from './user-interface-login-form.vue';
import UserInterfaceNotificationsForm from './user-interface-notifications-form.vue';

export default {
  components: {
    UserInterfaceGeneralForm,
    UserInterfaceLoginForm,
    UserInterfaceNotificationsForm,
    UserInterfaceAdvancedForm,
  },
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();
    const activeTab = ref('general');

    const generalElement = ref(null);
    const notificationsElement = ref(null);
    const advancedElement = ref(null);
    const loginPageElement = ref(null);

    const { hasChildrenError: hasGeneralError } = useValidationElementChildren(generalElement);
    const { hasChildrenError: hasNotificationsError } =
      useValidationElementChildren(notificationsElement);
    const { hasChildrenError: hasAdvancedError } = useValidationElementChildren(advancedElement);
    const { hasChildrenError: hasLoginPageError } = useValidationElementChildren(loginPageElement);

    return {
      t,
      activeTab,
      generalElement,
      notificationsElement,
      advancedElement,
      loginPageElement,
      hasGeneralError,
      hasNotificationsError,
      hasAdvancedError,
      hasLoginPageError,
    };
  },
};
</script>
