<template lang="pug">
  div.v-chip-group
    v-item-group(v-model="selected", :multiple="multiple")
      v-subheader(v-show="label") {{ label }}
      v-item(v-for="item in items", :key="item[itemValue]")
        v-chip(
          slot-scope="{ active, toggle }",
          :selected="active",
          :outline="outline",
          @click="toggle"
        ) {{ item[itemText] }}
</template>

<script>
import { isUndefined } from 'lodash';

import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'change',
  },
  props: {
    value: {
      type: [Array, Number, String],
      default() {
        return this.multiple ? [] : undefined;
      },
    },
    items: {
      type: Array,
      default: () => [],
    },
    itemText: {
      type: String,
      default: 'text',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    label: {
      type: String,
      default: null,
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    outline: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    selected: {
      get() {
        if (isUndefined(this.value)) {
          return this.value;
        }

        return this.multiple
          ? this.value.map(value => this.items.findIndex(item => item[this.itemValue] === value))
          : this.items.findIndex(item => item[this.itemValue] === this.value);
      },
      set(selected) {
        if (isUndefined(selected)) {
          return this.updateModel(selected);
        }

        const newModel = this.multiple
          ? selected.map(index => this.items[index][this.itemValue])
          : this.items[selected][this.itemValue];

        return this.updateModel(newModel);
      },
    },
  },
};
</script>

<style lang="scss" scoped>
  .v-chip-group ::v-deep {
    .v-subheader {
      height: 20px;
      font-size: 12px;
      font-weight: 400;
      padding: 0;
    }

    .v-item-group .v-chip {
      z-index: inherit;

      &, & .v-chip__content {
        z-index: inherit;
        cursor: pointer;
      }

      &:focus:not(.v-chip--disabled):not(.v-item--active):not(.v-chip--selected) {
        box-shadow: none;

        &:after {
          display: none;
        }
      }
    }
  }
</style>
