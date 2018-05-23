<template lang="pug">
  div.popups(v-if="popups.length")
    popup(v-bind="popups[0]", :key="popups[0].id")
    popup(v-if="popups[1]", v-bind="popups[1]", :key="popups[1].id")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import Popup from './popup.vue';

const { mapGetters } = createNamespacedHelpers('popup');

export default {
  components: { Popup },
  computed: {
    ...mapGetters(['popups']),
  },
};
</script>

<style lang="scss">
  .popups {
    position: fixed;
    z-index: 2;
    right: 2rem;
    top: 75px;
    width: 350px;
  }
</style>
