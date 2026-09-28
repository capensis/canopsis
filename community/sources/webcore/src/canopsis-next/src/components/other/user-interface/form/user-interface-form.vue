<template>
  <v-layout column>
    <v-tabs
      v-model="activeTab"
      slider-color="primary"
      centered
    >
      <v-tab :class="{ 'error--text': hasGeneralError }">
        {{ $t('common.general') }}
      </v-tab>
      <v-tab :class="{ 'error--text': hasNotificationsError }">
        {{ $tc('common.notification', 2) }}
      </v-tab>
      <v-tab :class="{ 'error--text': hasAdvancedError }">
        {{ $t('userInterface.advanced') }}
      </v-tab>
      <v-tab :class="{ 'error--text': hasLoginPageError }">
        {{ $t('userInterface.loginPage') }}
      </v-tab>

      <v-tab-item
        class="pt-4"
        eager
      >
        <user-interface-general-form
          v-field="form"
          ref="generalElement"
          :disabled="disabled"
        />
      </v-tab-item>

      <v-tab-item
        class="pt-4"
        eager
      >
        <user-interface-notifications-form
          v-field="form"
          ref="notificationsElement"
          :disabled="disabled"
        />
      </v-tab-item>

      <v-tab-item
        class="pt-4"
        eager
      >
        <user-interface-advanced-form
          v-field="form"
          ref="advancedElement"
          :disabled="disabled"
        />
      </v-tab-item>

      <v-tab-item
        class="pt-4"
        eager
      >
        <user-interface-login-form
          v-field="form"
          ref="loginPageElement"
          :disabled="disabled"
        />
      </v-tab-item>
    </v-tabs>
  </v-layout>
</template>

<script>
import { ref } from 'vue';

import { useValidationElementChildren } from '@/hooks/validator/validation-element-children';

import UserInterfaceAdvancedForm from './user-interface-advanced-form.vue';
import UserInterfaceGeneralForm from './user-interface-general-form.vue';
import UserInterfaceLoginForm from './user-interface-login-form.vue';
import UserInterfaceNotificationsForm from './user-interface-notifications-form.vue';

export default {
  inject: ['$validator'],
  components: {
    UserInterfaceGeneralForm,
    UserInterfaceLoginForm,
    UserInterfaceNotificationsForm,
    UserInterfaceAdvancedForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const activeTab = ref(0);

    const generalElement = ref(null);
    const loginPageElement = ref(null);
    const notificationsElement = ref(null);
    const advancedElement = ref(null);

    const { hasChildrenError: hasGeneralError } = useValidationElementChildren(generalElement);
    const { hasChildrenError: hasLoginPageError } = useValidationElementChildren(loginPageElement);
    const { hasChildrenError: hasNotificationsError } = useValidationElementChildren(notificationsElement);
    const { hasChildrenError: hasAdvancedError } = useValidationElementChildren(advancedElement);

    const reset = () => generalElement.value?.reset?.();

    return {
      activeTab,
      generalElement,
      loginPageElement,
      notificationsElement,
      advancedElement,
      hasGeneralError,
      hasLoginPageError,
      hasNotificationsError,
      hasAdvancedError,
      reset,
    };
  },
};
</script>
