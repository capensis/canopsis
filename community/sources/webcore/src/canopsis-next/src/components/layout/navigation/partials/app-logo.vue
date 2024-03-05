<template>
  <img
    :src="appLogo"
    alt=""
  >
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import defaultLogo from '@/assets/canopsis.png';

const { mapGetters } = createNamespacedHelpers('info');

export default {
  computed: {
    ...mapGetters(['logo']),

    appLogo() {
      if (this.logo) {
        return this.logo;
      }

      return defaultLogo;
    },
  },
};
</script>
