<template lang="pug">
  v-app#app
    v-layout(
      v-resize='onResize'
    )
      side-bar(
        :windowSize='windowSize'
      )
      top-bar
      v-content
        //
          v-btn(
          @click='handleClick'
          )
        router-view
</template>

<script>
import TopBar from '@/components/layout/top-bar.vue';
import SideBar from '@/components/layout/side-bar.vue';
import { mapState } from 'vuex';

export default {
  name: 'App',
  components: {
    TopBar,
    SideBar,
  },
  data() {
    return {
      windowSize: {
        x: 0,
        y: 0,
      },
    };
  },
  computed: {
    ...mapState({
      isSideBarOpen: state => state.app.isSideBarOpen,
    }),
  },
  methods: {
    handleClick() {
      // console.log(this.$store);
    },
    onResize() {
      this.windowSize = { x: window.innerWidth, y: window.innerHeight };
    },
  },
};
</script>

