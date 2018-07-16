<template lang="pug">
  v-toolbar.white(
    dense,
    fixed,
    clipped-left,
    app,
  )
    div.brand.ma-0.green.darken-4(v-show="$mq === 'tablet' || $mq === 'laptop'")
      img(src="../../assets/canopsis.png")
    v-toolbar-side-icon(@click="toggleSideBar")
    v-spacer
    div {{ $t('common.username') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('app');

/**
 * Component for the top bar of the application
 */
export default {
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
