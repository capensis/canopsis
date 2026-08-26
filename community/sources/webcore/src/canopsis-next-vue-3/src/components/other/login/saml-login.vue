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
const { samlConfig } = useInfo();

const title = computed(() => samlConfig.value?.title ?? t('login.loginWithSAML'));

const redirect = computed(() => route.query.redirect ?? '');

const samlHref = computed(() => {
  const loginUrl = removeTrailingSlashes(`${API_HOST}${API_ROUTES.saml.auth}`);
  const query = new URLSearchParams({
    relayState: removeTrailingSlashes(`${APP_HOST}${redirect.value}`),
  });

  return `${loginUrl}?${query}`;
});

const links = computed(() => [
  {
    title: title.value,
    href: samlHref.value,
  },
]);
</script>

<template>
  <third-party-login :title="t('login.loginWithSAML')" :links="links" />
</template>
