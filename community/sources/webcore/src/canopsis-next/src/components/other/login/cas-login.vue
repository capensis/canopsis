<template>
  <third-party-login
    :title="$t('login.loginWithCAS')"
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
      return get(this.casConfig, 'title', this.$t('login.loginWithCAS'));
    },

    casHref() {
      const { redirect = '' } = this.$route.query;

      const loginUrl = removeTrailingSlashes(`${API_HOST}${API_ROUTES.cas.login}`);
      const query = qs.stringify({
        redirect: removeTrailingSlashes(`${APP_HOST}${redirect}`),
        service: removeTrailingSlashes(`${API_HOST}${API_ROUTES.cas.loggedin}`),
      });

      return `${loginUrl}?${query}`;
    },

    links() {
      return [{
        title: this.title,
        href: this.casHref,
      }];
    },
  },
};
</script>
