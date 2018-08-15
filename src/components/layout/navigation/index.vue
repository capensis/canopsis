<template lang="pug">
  div
    top-bar(@toggleSideBar="toggleSideBar")
    side-bar(v-model="isOpened")
</template>

<script>
import TopBar from './top-bar.vue';
import SideBar from './side-bar.vue';

export default {
  components: {
    TopBar,
    SideBar,
  },
  data() {
    return {
      isOpened: false,
    };
  },
  methods: {
    toggleSideBar() {
      this.isOpened = !this.isOpened;
    },
  },
};
</script>
