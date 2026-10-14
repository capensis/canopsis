<template>
  <img :src="appLogo" alt="" />
</template>

<script>
import { useInfo } from '@/hooks/store/modules/info';

import defaultLogo from '@/assets/canopsis.png';

export default {
  setup() {
    const { logo } = useInfo();

    return {
      logo,
    };
  },

  computed: {
    appLogo() {
      return this.logo || defaultLogo;
    },
  },
};
</script>
