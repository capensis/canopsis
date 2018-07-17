<template lang="pug">
  v-navigation-drawer(
    v-model="isOpen",
    absolute,
    app,
    :clipped="$mq === 'mobile' || $mq === 'tablet' ? false : true"
  )
      v-card(flat)
      v-expansion-panel(
        class="panel",
        expand,
        focusable,
      )
        v-expansion-panel-content
          div(slot="header") Examples
          v-card
            v-card-text
              router-link(to="alarms") Alarms List
          v-card
            v-card-text
              router-link(to="filter") Filters
          v-card
            v-card-text
              router-link(to="login") Login
          v-card
            v-card-text
              router-link(to="rrule") Rrule
      v-divider
      v-expansion-panel(class="panel", expand, focusable)
        v-expansion-panel-content
          div(slot="header") View Group 2
          v-card
            v-card-text View 1
          v-card
            v-card-text View 2
      v-divider
      v-btn.addBtn(
        fab,
        dark,
        fixed,
        bottom,
        right,
        color="blue darken-4"
      )
        v-icon(dark) add
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import VueContentLoading from 'vue-content-loading';

const { mapGetters, mapActions } = createNamespacedHelpers('app');

/**
 * Component for the side-bar, on the left of the application
 */
export default {
  components: {
    VueContentLoading,
  },
  computed: {
    ...mapGetters(['isSideBarOpen']),

    isOpen: {
      get() {
        return this.isSideBarOpen;
      },
      set(state) {
        if (state !== this.isSideBarOpen) {
          this.toggleSideBar();
        }
      },
    },
  },
  methods: {
    ...mapActions(['toggleSideBar']),
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
