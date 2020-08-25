<template lang="pug">
  v-card
    v-card-text
      div.pa-3
        div.ml-2.mb-2.font-weight-bold {{ $t('login.loginWithCAS') }}
        v-btn.my-4(:href="casHref", color="primary") {{ title }}
</template>

<script>
import { get } from 'lodash';

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

      const { href: redirect } = this.$router.resolve({ name: 'home' });
      const { href } = this.$router.resolve({
        path: 'login',
        query: {
          service: `${this.casConfig.service}/logged_in`,
          redirect,
        },
      });

      return `${this.casConfig.server}${href}`;
    },
  },
};
</script>
