<template lang="pug">
  ul
    li.header.sticky(ref="header")
      v-checkbox.checkbox(
      v-if="checkbox",
      v-model="allSelected",
      hide-details
      )
      .headerText
        slot(name="header")
    transition(name="fade", mode="out-in")
      slot(name="loader", v-if="pending")
      div(v-else)
        li(v-for="item in items", :item="item")
          list-item(:item="item", :expanded="expanded")
            v-checkbox.checkbox(
            v-if="checkbox",
            slot="checkbox",
            :input-value="selected",
            :value="item._id",
            @change="$emit('update:selected',$event)",
            @click.stop,
            hide-details
            )
            .reduced(slot="reduced")
              slot(name="row", :props="item")
            div(slot="expanded")
              slot(name="expandedRow", :props="item")
        li(v-if="!items.length")
          div.container
            strong {{ $t('common.noResults') }}
</template>

<script>
import intersectionWith from 'lodash/intersectionWith';
import StickyFill from 'stickyfilljs';

import ListItem from '@/components/tables/list-item.vue';

/**
* Wrapper for lists (alarm list, entities list, ...)
*
* @prop {Array} [items] - Items to show on the list
* @prop {Boolean} [checkbox] - Boolean to determine if we need checkboxes on the table
* @prop {Boolean} [pending] - Boolean to know if loading is over
* @prop {Boolean} [expanded] - Boolean to know if rows are expanded or not
* @prop {Array} [selected] - Array of selected items (checkboxes)
*
* @event selected#update
*/
export default {
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
    pending: {
      type: Boolean,
      default: false,
    },
    expanded: {
      type: Boolean,
      default: false,
    },
    selected: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    allSelected: {
      get() {
        return this.selected.length === this.items.length && this.items.length !== 0;
      },
      set(value) {
        this.$emit('update:selected', value ? this.items.map(v => v._id) : []);
      },
    },
  },
  watch: {
    items(items) {
      const selected = intersectionWith(
        this.selected,
        items,
        (selectedItemId, item) => selectedItemId === item._id,
      );

      this.$emit('update:selected', selected);
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
    margin-left: 1%;
    &:hover {
      cursor: pointer;
    }
  }

  .fade-enter-active, .fade-leave-active {
    transition: opacity .5s;
  }

  .fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */
  {
    opacity: 0;
  }

  .headerText {
    margin-left: 1%;
  }
  .checkbox {
    position: absolute;
    top: 25%;
    width: 30px;
  }
</style>
