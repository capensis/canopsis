<template lang="pug">
    ul
      li.header.sticky(ref="header")
        v-checkbox.checkbox(v-if="checkbox",
        v-model="allSelected", hide-details, )
        .headerText
          slot(name="header")
      li(v-for="item in items", :item="item")
        list-item(:item="item")
          v-checkbox.checkbox(v-if="checkbox",
          v-model="selected", @change="$emit('update:selected',$event)",
          :value="item._id", @click.stop, hide-details, slot="checkbox")
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
      required: true,
    },
    checkbox: {
      type: Boolean,
      default: true,
    },
  },
  data() {
    return {
      // alarm's ids selected by the checkboxes
      selected: [],
    };
  },
  computed: {
    allSelected: {
      get() {
        return this.selected.length === this.items.length;
      },
      set(value) {
        this.selected = value ? this.items.map(v => v._id) : [];
        this.$emit('update:selected', this.selected);
      },
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
    top: 48px;
    z-index: 2;
  }
  .header {
    margin-bottom: 5px;
    background-color: rgb(251,247,247);
    z-index: 1;
  }
  .reduced {
    overflow: auto;
    margin-left: 1%;
  }
  .headerText {
    margin-left: 1%;
  }
  .checkbox {
    position: absolute;
    top: 25%;
  }
</style>
