<template lang="pug">
  v-navigation-drawer.side-bar.grey.darken-4(
    v-model="isOpen",
    absolute,
    app,
    :clipped="$mq === 'mobile' || $mq === 'tablet' ? false : true",
    :width="width"
  )
      v-card(flat)
      v-expansion-panel(
        class="panel",
        expand,
        focusable,
        dark
      )
        v-expansion-panel-content.grey.darken-4.white--text
          div(slot="header") Examples
          v-card.grey.darken-3.white--text
            v-card-text
              router-link(to="alarms") Alarms List
          v-card.grey.darken-3.white--text
            v-card-text
              router-link(to="context") Context
          v-card.grey.darken-3.white--text
            v-card-text
              router-link(to="weather") Service Weather
          v-card.grey.darken-3.white--text
            v-card-text
              router-link(to="filter") Filters
          v-card.grey.darken-3.white--text
            v-card-text
              router-link(to="login") Login
          v-card.grey.darken-3.white--text
            v-card-text
              router-link(to="rrule") Rrule
      v-expansion-panel(
        class="panel",
        expand,
        focusable,
        dark
      )
        v-expansion-panel-content.grey.darken-4.white--text
          div(slot="header") View Group 2
          v-card.grey.darken-3.white--text
            v-card-text View 1
          v-card.grey.darken-3.white--text
            v-card-text View 2
      v-divider
      v-btn.addBtn(
        fab,
        dark,
        fixed,
        bottom,
        right,
        color="green darken-4"
      )
        v-icon(dark) add
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import VueContentLoading from 'vue-content-loading';

import { SIDE_BAR_WIDTH } from '@/config';

const { mapGetters, mapActions } = createNamespacedHelpers('app');

/**
 * Component for the side-bar, on the left of the application
 */
export default {
  components: {
    VueContentLoading,
  },
  data() {
    return {
      width: SIDE_BAR_WIDTH,
    };
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

  .side-bar {
    position: fixed;
    height: 100vh;
    overflow-y: scroll;
  }
</style>
