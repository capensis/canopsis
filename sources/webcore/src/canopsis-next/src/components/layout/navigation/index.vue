<template lang="pug">
  div
    top-bar(@toggleSideBar="toggleSideBar")
    groups-side-bar(v-if="isShownGroupsSideBar", v-model="isOpen")
</template>

<script>
import appMixin from '@/mixins/app';
import entitiesInfoMixin from '@/mixins/entities/info';

import TopBar from './partial/top-bar.vue';
import GroupsSideBar from './partial/groups-side-bar.vue';

export default {
  components: {
    TopBar,
    GroupsSideBar,
  },
  mixins: [appMixin, entitiesInfoMixin],
  data() {
    return {
      isOpen: false,
    };
  },
  mounted() {
    this.fetchAppInfos();
  },
  methods: {
    toggleSideBar() {
      this.isOpen = !this.isOpen;
    },
  },
};
</script>
