<template lang="pug">
  div
    top-bar(@toggleSideBar="toggleSideBar")
    groups-side-bar(v-if="isShownGroupsSideBar", v-model="isOpen")
</template>

<script>
import authMixin from '@/mixins/auth';

import TopBar from '@/components/common/navigation/top-bar.vue';
import GroupsSideBar from '@/components/common/navigation/groups-side-bar.vue';

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
