<template lang="pug">
  div
    top-bar(@toggleSideBar="toggleSideBar")
    side-bar(v-model="isOpen")
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
