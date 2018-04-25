<template lang="pug">
  v-navigation-drawer(
    v-model='isOpen'
    absolute
    app
    :clipped="windowSize.x < 1264 ? false : true"
  )
      v-card(
        flat
      )
      v-expansion-panel(
        class='panel'
        expand
        focusable
      )
        v-expansion-panel-content
          div(slot="header") View Group 1
          v-card
            v-card-text View 1
          v-card
            v-card-text View 2
      v-divider
      v-expansion-panel(
        class='panel'
        expand
        focusable
      )
        v-expansion-panel-content
          div(slot="header") View Group 2
          v-card
            v-card-text View 1
          v-card
            v-card-text View 2
      v-divider
        v-btn(
          class='addBtn'
          fab
          dark
          fixed
          bottom
          right
          color="blue darken-4"
        )
          v-icon(dark) add
</template>

<script>
import VueContentLoading from 'vue-content-loading';

export default {
  name: 'SideBar',
  components: {
    VueContentLoading,
  },
  props: {
    windowSize: {
      type: Object,
    },
  },
  computed: {
    isOpen: {
      get() {
        return this.$store.state.app.isSideBarOpen;
      },
      set(state) {
        if (state !== this.$store.state.app.isSideBarOpen) {
          this.$store.dispatch('app/toggleSideBar');
        }
      },
    },
  },
};
</script>

<style scoped>
  a {
    color: inherit;
    text-decoration: none;
  }
  .panel {
    box-shadow: none;
  }
</style>
