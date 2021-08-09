<template lang="pug">
  v-layout.error-page__wrapper(align-center)
    v-layout(column, align-center)
      img(src="@/assets/canopsis-green.png")
      span.subheading.pt-4 {{ $t('errors.default') }}
      span.subheading(v-if="message") {{ message }}
      v-btn.primary.mt-4(small, @click="redirectToLogin") {{ $t('common.login') }}
</template>

<script>
import { ROUTE_NAMES } from '@/constants';

export default {
  props: {
    message: {
      type: String,
      default: '',
    },
  },
  methods: {
    redirectToLogin() {
      return this.$router.push({ name: ROUTE_NAMES.login });
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
