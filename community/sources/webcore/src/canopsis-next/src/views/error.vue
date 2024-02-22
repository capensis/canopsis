<template>
  <v-layout
    class="error-page__wrapper"
    align-center
  >
    <v-layout
      column
      align-center
    >
      <img
        src="@/assets/canopsis-green.png"
        alt=""
      >
      <span class="text-subtitle-1 pt-4">{{ $t('errors.default') }}</span>
      <span
        v-if="message"
        class="text-subtitle-1"
      >
        {{ message }}
      </span>
    </v-layout>
  </v-layout>
</template>

<script>
import { isEmpty } from 'lodash';

import { APP_INFO_FETCHING_INTERVAL } from '@/config';
import { ROUTES_NAMES } from '@/constants';

import { authMixin } from '@/mixins/auth';
import { entitiesInfoMixin } from '@/mixins/entities/info';
import { pollingMixinCreator } from '@/mixins/polling';

export default {
  mixins: [
    authMixin,
    entitiesInfoMixin,
    pollingMixinCreator({
      method: 'fetchInfo',
      delayField: 'pollingDelay',
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
  computed: {
    pollingDelay() {
      return APP_INFO_FETCHING_INTERVAL;
    },
  },
  mounted() {
    this.fetchInfo();
  },
  methods: {
    async fetchInfo() {
      try {
        await this.fetchAppInfo();

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
