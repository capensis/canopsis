<script setup>
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute, useRouter } from 'vue-router';

import { ROUTES_NAMES, ROUTES, RESPONSE_STATUSES } from '@/constants';
import { useAuth } from '@/hooks/auth';
import { useForm } from '@/hooks/form';
import { useInfo } from '@/hooks/store/modules/info';

import LdapLoginInformation from './ldap-login-information.vue';
import LoginForm from './login-form.vue';

const { currentUser, login } = useAuth();
const { isLDAPAuthEnabled, footer } = useInfo();

const route = useRoute();
const router = useRouter();
const { t } = useI18n();

const serverErrorMessage = ref('');

const form = useForm({
  initialValues: {
    username: '',
    password: '',
  },
});
const { handleSubmit, isSubmitting } = form;

const submit = handleSubmit(async () => {
  try {
    serverErrorMessage.value = '';

    const { redirect } = route.query;

    await login(form.values.value);

    const { defaultview: userDefaultView } = currentUser.value ?? {};

    if (redirect && redirect !== ROUTES.home) {
      await router.push(redirect);
      return;
    }

    if (userDefaultView) {
      await router.push({
        name: ROUTES_NAMES.view,
        params: { id: userDefaultView._id },
      });
      return;
    }

    await router.push({ name: ROUTES_NAMES.home });
  } catch (err) {
    serverErrorMessage.value =
      err?.status === RESPONSE_STATUSES.serviceUnavailable
        ? t('login.errors.underMaintenance')
        : t('login.errors.incorrectEmailOrPassword');
  }
});
</script>

<template>
  <v-form @submit.prevent.stop="submit">
    <c-layout column>
      <div>
        <ldap-login-information v-if="isLDAPAuthEnabled" />
      </div>
      <login-form />

      <div>
        <div class="d-flex justify-space-between align-center">
          <v-btn :loading="isSubmitting" type="submit" color="primary">
            {{ t('common.connect') }}
          </v-btn>

          <div class="flex-grow-1 ml-4">
            <c-alert v-if="serverErrorMessage" class="py-1 my-0 font-weight-bold" type="error">
              {{ serverErrorMessage }}
            </c-alert>
          </div>
        </div>

        <template v-if="footer">
          <v-divider class="my-2" />

          <div>
            <c-compiled-template :template="footer" />
          </div>
        </template>
      </div>
    </c-layout>
  </v-form>
</template>
