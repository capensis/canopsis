<template lang="pug">
  div.white.treeview-data-table
    v-layout(row)
      div.treeview-data-table--tree.mr-4
        v-treeview(
          :open.sync="opened",
          :items="items",
          :item-children="itemChildren",
          :item-key="itemKey",
          :open-all="openAll",
          :load-children="loadChildren",
          :dark="dark",
          :light="light"
        )
          template(slot="label", slot-scope="{ item }")
            slot(name="expand", :item="item")
              v-avatar.white--text(color="primary", size="32") {{ item | get(`${itemChildren}.length`, null, 0) }}
            slot(name="expand-append", :item="item")
      v-flex
        v-data-table(
          :headers="headers",
          :items="openedItems",
          :loading="loading",
          :total-items="totalItems",
          :header-text="headerText",
          :item-key="itemKey",
          :dark="dark",
          :light="light",
          hide-actions
        )
          template(slot="items", slot-scope="props")
            slot(v-bind="props", name="items")
              tr(:key="props.item[itemKey]")
                td(v-for="header in headers", :key="header.value")
                  slot(:name="header.value", v-bind="props") {{ props.item | get(header.value) }}
</template>

<script>
import { convertTreeArrayToArray } from '@/helpers/treeview';

export default {
  props: {
    headers: {
      type: Array,
      required: true,
    },
    items: {
      type: Array,
      required: true,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    openAll: {
      type: Boolean,
      default: false,
    },
    loadChildren: {
      type: Function,
      default: null,
    },
    itemKey: {
      type: String,
      default: '_id',
    },
    itemChildren: {
      type: String,
      default: 'children',
    },
    headerText: {
      type: String,
      default: 'text',
    },
    dark: {
      type: Boolean,
      default: false,
    },
    light: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      opened: [],
    };
  },
  computed: {
    totalItems() {
      return this.items.length;
    },

    openedItems() {
      const arrayItems = convertTreeArrayToArray(this.items, this.itemChildren);
      const itemsById = arrayItems.reduce((acc, item) => {
        acc[item[this.itemKey]] = item;

        return acc;
      }, {});

      const isItemOpen = (item = {}) =>
        !item.parentKey || (this.opened.includes(item.parentKey) && isItemOpen(itemsById[item.parentKey]));

      return arrayItems.filter(isItemOpen);
    },
  },
};
</script>

<style lang="scss" scoped>
.treeview-data-table {
  /deep/ .v-treeview-node {
    margin-left: 34px;

    &--leaf {
      margin-left: 58px;
    }

    &__root {
      height: 48px;
    }

    &__label {
      .v-btn .v-icon, .v-avatar .v-icon {
        padding-right: 0;
      }
    }
  }

  /deep/ .v-treeview > .v-treeview-node {
    margin-left: 0;
  }

  &--tree {
    margin-top: 56px;
  }
}
</style>