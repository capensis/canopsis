<template lang="pug">
  div.popups(v-if="popups.length")
    popup(v-for="popup in popups", v-bind="popup", :key="popup.id")
</template>

<script>
// LIBS
import { createNamespacedHelpers } from 'vuex';
// COMPONENTS
import Popup from '@/components/popups/popup.vue';

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
