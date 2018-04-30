<template lang="pug">
    ul
      li.header.sticky(ref="header" :style="headerStyle")
        slot(name="header")
      li(v-for="item in items", :item="item")
        brick-list
          div(slot="reduced", :style="reducedStyle")
            slot(name="row", :props="item")
          div(slot="expanded")
            slot(name="expandedRow", :props="item")
</template>

<script>
import StickyFill from 'stickyfilljs';
import BrickList from './brick-list.vue';

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
        backgroundColor: 'rgb(251,247,247)  ',
        zIndex: '1',
      },
    };
  },
  mounted() {
    StickyFill.addOne(this.$refs.header);
  },
  beforeDestroy() {
    StickyFill.removeOne(this.$refs.header);
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
  .sticky {
    position: -webkit-sticky;
    position: sticky;
    top: 48px;
    z-index: 2;
  }
</style>
