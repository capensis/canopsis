<template lang="pug">
  v-layout.error-page__wrapper(align-center)
    v-layout(column, align-center)
      img(src="@/assets/canopsis-green.png")
      span.subheading.pt-4 {{ $t('errors.default') }}
      span.subheading(v-if="message") {{ message }}
</template>

<script>
import { isEmpty } from 'lodash';

import { LOGIN_INFOS_FETCHING_INTERVAL } from '@/config';

import { ROUTES_NAMES } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { entitiesInfoMixin } from '@/mixins/entities/info';
import { pollingMixinCreator } from '@/mixins/polling';

export default {
  mixins: [
    authMixin,
    entitiesInfoMixin,
    pollingMixinCreator({
      method: 'fetchInfos',
      delay: LOGIN_INFOS_FETCHING_INTERVAL,
      startOnMount: true,
    }),
  ],
  props: {
    message: {
      type: String,
      default: '',
    },
    redirect: {
      type: String,
      default: '',
    },
  },
  mounted() {
    this.fetchInfos();
  },
  methods: {
    async fetchInfos() {
      try {
        await this.fetchLoginInfos();

        if (!isEmpty(this.currentUser)) {
          this.$router.replace({
            name: this.redirect || ROUTES_NAMES.home,
          });
        } else {
          this.$router.replace({
            name: ROUTES_NAMES.login,
            query: {
              redirect: this.redirect,
            },
          });
        }
      } catch (err) {
        console.error(err);
      }
    },
  },
};
</script>

<style lang="scss">
.error-page {
  &__wrapper {
    min-height: 100vh;
  }
}
</style>
