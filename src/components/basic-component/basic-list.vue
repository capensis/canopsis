<template lang="pug">
    ul
      li.header.sticky(ref="header")
        slot(name="header")
      transition(name="fade", mode="out-in")
        slot(name="loader", v-if="pending")
        div(v-else)
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
    pending: {
      type: Boolean,
      default: false,
    },
  },
  mounted() {
    const { header } = this.$refs;
    if (header) {
      StickyFill.addOne(header);
    }
  },
  beforeDestroy() {
    const { header } = this.$refs;
    if (header) {
      StickyFill.removeOne(header);
    }
  },
};
</script>

<style scoped lang="scss">
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
  .header {
    z-index: 1;
    font-size: 0.9em;
    line-height: 2em;
    background-color: white;
  }
  .reduced {
    overflow: auto;

    &:hover {
      cursor: pointer;
    }
  }
  .fade-enter-active, .fade-leave-active {
    transition: opacity .5s;
  }
  .fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
    opacity: 0;
  }
</style>
