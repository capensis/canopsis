<template lang="pug">
  v-sheet.treeview-data-table
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
          template(#label="{ item }")
            slot(name="expand", :item="item")
              v-avatar.white--text(color="primary", size="32") {{ item | get(`${itemChildren}.length`, 0) }}
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
          template(#items="props")
            slot(v-bind="props", name="items")
              tr(:key="props.item[itemKey]")
                td(v-for="header in headers", :key="header.value")
                  slot(:name="header.value", v-bind="props") {{ props.item | get(header.value) }}
</template>

<script>
import { keyBy } from 'lodash';

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

    arrayItems() {
      return convertTreeArrayToArray(this.items, this.itemChildren);
    },

    itemsById() {
      return keyBy(this.arrayItems, this.itemKey);
    },

    openedItems() {
      return this.arrayItems.filter(item => this.isDependencyOpen(item));
    },
  },
  methods: {
    isDependencyOpen(item, parentKeys = []) {
      if (!item.parentKey || parentKeys.includes(item.key)) {
        return true;
      }

      return (
        this.opened.includes(item.parentKey)
        && this.isDependencyOpen(this.itemsById[item.parentKey], [...parentKeys, item.key])
      );
    },
  },
};
</script>

<style lang="scss" scoped>
.treeview-data-table {
  ::v-deep .v-treeview-node {
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

  ::v-deep .v-treeview > .v-treeview-node {
    margin-left: 0;
  }

  &--tree {
    margin-top: 56px;
  }
}
</style>
