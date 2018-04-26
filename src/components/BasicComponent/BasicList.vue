<template lang="pug">
    ul
      li(ref="header" class='header')
        sticky-wrapper()
          slot(name="header")
      li(v-for="item in items", :item="item")
        brick-list
          div(slot="reduced", :style="reducedStyle")
            slot(name="row", :props="item")
          div(slot="expanded")
            slot(name="expandedRow", :props="item")
</template>

<script>
import StickyWrapper from './StickyWrapper.vue';
import BrickList from './BrickList.vue';

export default {
  name: 'BasicList',
  components: { BrickList, StickyWrapper },
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
  methods: {
    changePage() {
      // On fait appel au store pour charger les nouvelles entités
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
  }
  .header {
    border-bottom: 1px solid gray;
    background-color: white;
    position: sticky;
    top: 50px;
    z-index: 500;
  }
</style>
