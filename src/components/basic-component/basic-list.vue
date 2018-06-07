<template lang="pug">
    ul
      li.header.sticky(ref="header")
        slot(name="header")
      li(v-for="item in items", :item="item")
        list-item
          .reduced(slot="reduced")
            slot(name="row", :props="item")
          div(slot="expanded")
            slot(name="expandedRow", :props="item")
      li(v-if="!items.length")
        div.container
          strong {{ $t('common.noResults') }}
</template>

<script>
import StickyFill from 'stickyfilljs';
import ListItem from '@/components/basic-component/list-item.vue';

export default {
  name: 'BasicList',
  components: { ListItem },
  props: {
    items: {
      type: Array,
    },
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
  ul {
   position: relative;
   list-style-type: none;
  }
  .sticky {
    position: -webkit-sticky;
    position: sticky;
    top: 38px;
    z-index: 2;
  }
  .header {
    background-color: rgb(251,247,247);
    z-index: 1;
  }
  .reduced {
    overflow: auto;
  }
</style>
