<script setup>
import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useRoute } from 'vue-router';

import { API_HOST, API_ROUTES, APP_HOST } from '@/config';

import { removeTrailingSlashes } from '@/helpers/url';

import { useInfo } from '@/hooks/store/modules/info';

import ThirdPartyLogin from './third-party-login.vue';

const route = useRoute();
const { t } = useI18n();
const { oauthConfig } = useInfo();

const redirect = computed(() => route.query.redirect ?? '');

const links = computed(() =>
  (oauthConfig.value?.providers ?? []).map((provider) => {
    const loginUrl = removeTrailingSlashes(`${API_HOST}${API_ROUTES.oauth.login}${provider}/login`);
    const query = new URLSearchParams({
      redirect: removeTrailingSlashes(`${APP_HOST}${redirect.value}`),
    });

    return {
      title: `${t('common.login')} ${provider}`,
      href: `${loginUrl}?${query}`,
    };
  }),
);
</script>

<template>
  <third-party-login :title="t('login.loginWithOauth')" :links="links" />
</template>
