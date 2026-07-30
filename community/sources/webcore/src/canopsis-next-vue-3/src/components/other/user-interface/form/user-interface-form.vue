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
        <user-interface-general-form :disabled="disabled" />
      </v-window-item>

      <v-window-item value="notifications" class="pt-4" eager>
        <user-interface-notifications-form :disabled="disabled" />
      </v-window-item>

      <v-window-item value="advanced" class="pt-4" eager>
        <user-interface-advanced-form :disabled="disabled" />
      </v-window-item>

      <v-window-item value="login" class="pt-4" eager>
        <user-interface-login-form :disabled="disabled" />
      </v-window-item>
    </v-window>
  </c-layout>
</template>

<script>
import { computed, ref } from 'vue';

import { useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import UserInterfaceAdvancedForm from './user-interface-advanced-form.vue';
import UserInterfaceGeneralForm from './user-interface-general-form.vue';
import UserInterfaceLoginForm from './user-interface-login-form.vue';
import UserInterfaceNotificationsForm from './user-interface-notifications-form.vue';

const hasErrorByPrefixes = (errors, prefixes) =>
  Object.keys(errors).some(name =>
    prefixes.some(prefix => name === prefix || name.startsWith(`${prefix}.`)),
  );

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
    const form = useFormContext();
    const activeTab = ref('general');

    const hasGeneralError = computed(() =>
      hasErrorByPrefixes(form.errors.value, [
        'app_title',
        'language',
        'default_color_theme',
        'timezone',
        'logo',
        'version_description',
      ]),
    );
    const hasNotificationsError = computed(() =>
      hasErrorByPrefixes(form.errors.value, ['instruction']),
    );
    const hasAdvancedError = computed(() =>
      hasErrorByPrefixes(form.errors.value, [
        'popup_timeout',
        'max_matched_items',
        'check_count_request_timeout',
        'show_header_on_kiosk_mode',
        'disabled_transitions',
        'allow_change_severity_to_info',
        'required_instruction_approve',
        'auto_suggest_pbehavior_name',
      ]),
    );
    const hasLoginPageError = computed(() =>
      hasErrorByPrefixes(form.errors.value, ['footer', 'login_page_description']),
    );

    return {
      t,
      activeTab,
      hasGeneralError,
      hasNotificationsError,
      hasAdvancedError,
      hasLoginPageError,
    };
  },
};
</script>
