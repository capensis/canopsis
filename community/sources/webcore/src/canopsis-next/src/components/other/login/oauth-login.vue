<template>
  <third-party-login
    :title="$t('login.loginWithOauth')"
    :links="links"
  />
</template>

<script>
import qs from 'qs';

import { API_HOST, API_ROUTES, APP_HOST } from '@/config';

import { removeTrailingSlashes } from '@/helpers/url';

import { entitiesInfoMixin } from '@/mixins/entities/info';

import ThirdPartyLogin from './partials/third-party-login.vue';

export default {
  components: { ThirdPartyLogin },
  mixins: [entitiesInfoMixin],
  computed: {
    links() {
      return (this.oauthConfig?.providers ?? []).map((provider) => {
        const { redirect = '' } = this.$route.query;

        const loginUrl = removeTrailingSlashes(`${API_HOST}${API_ROUTES.oauth.login}${provider}/login`);
        const query = qs.stringify({
          redirect: removeTrailingSlashes(`${APP_HOST}${redirect}`),
        });

        return {
          title: `${this.$t('common.login')} ${provider}`,
          href: `${loginUrl}?${query}`,
        };
      });
    },
  },
};
</script>
