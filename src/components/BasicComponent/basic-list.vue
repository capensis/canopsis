<template lang="pug">
    ul
      li.header.sticky(ref="header", :style="headerStyle")
        slot(name="header")
      li(v-for="item in items", :item="item")
        list-item
          div(slot="reduced", :style="reducedStyle")
            slot(name="row", :props="item")
          div(slot="expanded")
            slot(name="expandedRow", :props="item")
</template>

<script>
import StickyFill from 'stickyfilljs';
import ListItem from '@/components/BasicComponent/list-item.vue';

export default {
  name: 'BasicList',
  components: { ListItem },
  props: {
    items: {
      type: Array,
    },
  },
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
