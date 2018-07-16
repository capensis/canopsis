<template lang="pug">
  v-toolbar.white(
    dense,
    fixed,
    clipped-left,
    app,
  )
    div.brand.ma-0.green.darken-4(v-show="$mq === 'tablet' || $mq === 'desktop'")
      img(:src="canopsisLogo")
    v-toolbar-side-icon(@click="toggleSideBar")
    v-spacer
    div {{ $t('common.username') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import canopsisLogo from '@/assets/canopsis.png';

const { mapActions } = createNamespacedHelpers('app');

/**
 * Component for the top bar of the application
 */
export default {
  data() {
    return {
      canopsisLogo,
    };
  },
  methods: {
    ...mapActions([
      'toggleSideBar',
    ]),
  },
};
</script>

<style lang="scss" scoped>
  .brand {
    display: flex;
    align-items: center;
    margin: 0;
    width: 250px;
    height: 100%;

    img {
      margin: auto;
    }
  }
</style>
