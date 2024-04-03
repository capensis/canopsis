<template>
  <third-party-login
    :title="$t('login.loginWithSAML')"
    :links="links"
  />
</template>

<script>
import qs from 'qs';
import { get } from 'lodash';

import { APP_HOST, API_HOST, API_ROUTES } from '@/config';

import { removeTrailingSlashes } from '@/helpers/url';

import { entitiesInfoMixin } from '@/mixins/entities/info';

import ThirdPartyLogin from './partials/third-party-login.vue';

export default {
  components: { ThirdPartyLogin },
  mixins: [entitiesInfoMixin],
  computed: {
    title() {
      return get(this.samlConfig, 'title', this.$t('login.loginWithSAML'));
    },

    samlHref() {
      const { redirect = '' } = this.$route.query;

      const loginUrl = removeTrailingSlashes(`${API_HOST}${API_ROUTES.saml.auth}`);
      const query = qs.stringify({ relayState: removeTrailingSlashes(`${APP_HOST}${redirect}`) });

      return `${loginUrl}?${query}`;
    },

    links() {
      return [{
        title: this.title,
        href: this.samlHref,
      }];
    },
  },
};
</script>
