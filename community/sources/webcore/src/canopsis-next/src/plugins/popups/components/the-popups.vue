<template>
  <div
    class="popups"
    v-if="popups.length"
  >
    <popup-item
      v-for="popup in popups"
      v-bind="popup"
      :key="popup.id"
    />
  </div>
</template>

<script>
import PopupItem from './popup-item.vue';

/**
 * Wrapper for the popups
 */
export default {
  components: { PopupItem },
  computed: {
    popups() {
      return this.$store.getters[`${this.$popups.moduleName}/popups`];
    },
  },
};
</script>

<style lang="scss">
  .popups {
    position: fixed;
    z-index: 1000;
    right: 2rem;
    top: 75px;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
  }
</style>
