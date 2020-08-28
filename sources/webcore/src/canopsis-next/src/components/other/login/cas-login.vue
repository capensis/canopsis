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

import { getApplicationHost } from '@/helpers/router';

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

      const { href } = this.$router.resolve({ name: 'home' });
      const { redirect = href } = this.$route.query;
      const appUrl = getApplicationHost();
      const query = qs.stringify({
        redirect: `${appUrl}${redirect}`,
        service: `${appUrl}api/cas/loggedin`,
      });

      return `${appUrl}api/cas/login?${query}`;
    },
  },
};
</script>
