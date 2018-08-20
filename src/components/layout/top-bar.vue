<template lang="pug">
  v-toolbar.white(
    dense,
    fixed,
    clipped-left,
    app,
  )
    div.brand.ma-0.green.darken-4(v-show="$mq | mq({ t: true })")
      img(src="../../assets/canopsis.png")
    v-toolbar-side-icon(@click="toggleSideBar")
    v-spacer
    v-toolbar-items
      v-menu(offset-y, bottom)
        v-btn(slot="activator", flat) {{ currentUser.crecord_name }}
        v-list
          v-list-tile(@click.prevent="logout")
            v-list-tile-title {{ $t('common.logout') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import authMixin from '@/mixins/auth';

const { mapActions } = createNamespacedHelpers('app');

/**
 * Component for the top bar of the application
 */
export default {
  mixins: [authMixin],
  methods: {
    ...mapActions(['toggleSideBar']),
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
