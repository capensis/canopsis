<template lang="pug">
  div.popups(v-if="popups.length", data-test="popupsWrapper")
    popup(v-for="popup in popups", v-bind="popup", :key="popup.id")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import Popup from './popup.vue';

const { mapGetters } = createNamespacedHelpers('popup');

/**
* Wrapper for the popups
*/
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
    z-index: 1000;
    right: 2rem;
    top: 75px;
  }
</style>
