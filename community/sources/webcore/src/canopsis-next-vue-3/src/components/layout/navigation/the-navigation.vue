<template>
  <div>
    <groups-side-bar v-if="isShownGroupsSideBar" v-model="isOpen" />
    <top-bar :side-bar-open="isOpen" @toggle:side-bar="toggleSideBar" />
  </div>
</template>

<script>
import { authMixin } from '@/mixins/auth';

import TopBar from './partials/top-bar.vue';
import GroupsSideBar from './partials/groups-side-bar/groups-side-bar.vue';

export default {
  name: 'TheNavigation',

  components: {
    GroupsSideBar,
    TopBar,
  },

  mixins: [authMixin],

  data() {
    return {
      isOpen: false,
    };
  },

  methods: {
    toggleSideBar() {
      this.isOpen = !this.isOpen;
    },
  },
};
</script>
