<template>
  <third-party-login :title="t('login.loginWithCAS')" :links="links" />
</template>

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
const { casConfig } = useInfo();

const title = computed(() => casConfig.value?.title ?? t('login.loginWithCAS'));

const redirect = computed(() => route.query.redirect ?? '');

const casHref = computed(() => {
  const loginUrl = removeTrailingSlashes(`${API_HOST}${API_ROUTES.cas.login}`);
  const query = new URLSearchParams({
    redirect: removeTrailingSlashes(`${APP_HOST}${redirect.value}`),
    service: removeTrailingSlashes(`${API_HOST}${API_ROUTES.cas.loggedin}`),
  });

  return `${loginUrl}?${query}`;
});

const links = computed(() => [
  {
    title: title.value,
    href: casHref.value,
  },
]);
</script>
