<template>
  <div>
    <top-bar @toggleSideBar="toggleSideBar" />
    <groups-side-bar
      v-if="isShownGroupsSideBar"
      v-model="isOpen"
    />
  </div>
</template>

<script>
import { authMixin } from '@/mixins/auth';

import TopBar from './partials/top-bar.vue';
import GroupsSideBar from './partials/groups-side-bar/groups-side-bar.vue';

export default {
  components: {
    TopBar,
    GroupsSideBar,
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
