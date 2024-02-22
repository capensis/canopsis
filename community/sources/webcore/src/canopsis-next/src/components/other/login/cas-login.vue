<template>
  <v-card>
    <v-card-text>
      <div class="ml-2 mb-2 font-weight-bold">
        {{ $t('login.loginWithCAS') }}
      </div>
      <v-btn
        :href="casHref"
        color="primary"
      >
        {{ title }}
      </v-btn>
    </v-card-text>
  </v-card>
</template>

<script>
import qs from 'qs';
import { get } from 'lodash';

import { APP_HOST, API_HOST, API_ROUTES } from '@/config';

import { removeTrailingSlashes } from '@/helpers/url';

import { entitiesInfoMixin } from '@/mixins/entities/info';

export default {
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
  },
};
</script>
