<template>
  <img
    :src="appLogo"
    alt=""
  >
</template>

<script>
import defaultLogo from '@/assets/canopsis.png';
import { useInfo } from '@/hooks/store/modules/info';

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
