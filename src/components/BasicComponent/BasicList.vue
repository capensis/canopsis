<template lang="pug">
    ul
        li(ref="header")
          slot(name="header")
        li(v-for="item in items" :item="item")
          brick-list
            div(slot="reduced" :style="reducedStyle")
              slot(name="row" :props="item")
            div(slot="expanded")
              slot(name="expandedRow" :props="item")
</template>

<script>
import BrickList from './BrickList.vue';

export default {
  name: 'BasicList',
  components: { BrickList },
  props: ['items'],
  data() {
    return {
      reducedStyle: {
        overflow: 'auto',
      },
      headerStyle: {
        marginBottom: '5px',
        zIndex: '1',
      },
    };
  },
  created() {
    document.addEventListener('scroll', this.stickyHeader);
  },
  destroyed() {
    document.removeEventListener('scroll', this.stickyHeader);
  },
  methods: {
    changePage() {
      // On fait appel au store pour charger les nouvelles entités
    },
    stickyHeader() {
      if (window.pageYOffset >= this.$refs.header.getBoundingClientRect().top) {
        this.$refs.header.classList.add('sticky');
      } else {
        this.$refs.header.classList.remove('sticky');
      }
    },
  },
};
</script>

<style scoped>
  * {
    box-sizing:border-box;
  }
  ul {
   position: relative;
   list-style-type: none;
   height: 100%;
}
  .sticky {
    position: fixed;
    top: 0;
    z-index: 1;
  }
</style>
