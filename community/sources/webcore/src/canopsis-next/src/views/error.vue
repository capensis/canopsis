<template lang="pug">
  v-layout.error-page__wrapper(align-center)
    v-layout(column, align-center)
      img(src="@/assets/canopsis-green.png")
      span.subheading.pt-4 {{ $t('errors.default') }}
      span.subheading(v-if="message") {{ message }}
</template>

<script>
import { LOGIN_INFOS_FETCHING_INTERVAL } from '@/config';

import { ROUTES_NAMES } from '@/constants';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { createPollingMixin } from '@/mixins/polling';

export default {
  mixins: [
    entitiesInfoMixin,
    createPollingMixin({
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
  },
  mounted() {
    this.fetchInfos();
  },
  methods: {
    async fetchInfos() {
      try {
        await this.fetchLoginInfos();

        this.$router.replace({ name: ROUTES_NAMES.login });
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
