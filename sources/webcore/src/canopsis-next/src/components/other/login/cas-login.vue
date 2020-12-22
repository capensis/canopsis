<template lang="pug">
  v-card
    v-card-text
      div.pa-3
        div.ml-2.mb-2.font-weight-bold {{ $t('login.loginWithCAS') }}
        v-btn.my-4(:href="casHref", color="primary") {{ title }}
</template>

<script>
import qs from 'qs';
import { get } from 'lodash';

import { APP_HOST, API_HOST, API_ROUTES } from '@/config';

import entitiesInfoMixin from '@/mixins/entities/info';

export default {
  mixins: [entitiesInfoMixin],
  computed: {
    title() {
      return get(this.casConfig, 'title', this.$t('login.loginWithCAS'));
    },

    casHref() {
      if (!this.casConfig) {
        return null;
      }

      const { redirect = '' } = this.$route.query;

      const query = qs.stringify({
        redirect: `${APP_HOST}${redirect}`,
        service: `${API_HOST}${API_ROUTES.cas.loggedin}`,
      });

      return `${API_HOST}${API_ROUTES.cas.login}?${query}`;
    },
  },
};
</script>
