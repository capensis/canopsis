<template lang="pug">
  component.c-select-field(
    :value="value",
    v-bind="$attrs",
    v-on="$listeners",
    v-validate="rules",
    :class="{ 'c-select-field--ellipsis': ellipsis }",
    :is="component",
    :item-text="itemText",
    :item-value="itemValue",
    ref="select"
  )
    template(v-if="$scopedSlots.selection || ellipsis", #selection="props")
      slot(name="selection", v-bind="props")
        span.ellipsis {{ props.item[itemText] }}
    template(v-if="$scopedSlots.item", #item="props")
      slot(name="item", v-bind="props")
    template(v-if="$scopedSlots['append-item']", #append-item="")
      slot(name="append-item")
</template>

<script>
import { isArray } from 'lodash';

export default {
  inject: ['$validator'],
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Array, Object, String, Symbol],
      default: '',
    },
    required: {
      type: Boolean,
      default: false,
    },
    autocomplete: {
      type: Boolean,
      default: false,
    },
    combobox: {
      type: Boolean,
      default: false,
    },
    ellipsis: {
      type: Boolean,
      default: false,
    },
    itemText: {
      type: [String, Function],
      default: 'text',
    },
    itemValue: {
      type: String,
      default: 'value',
    },
  },
  computed: {
    isArray() {
      return isArray(this.value);
    },

    component() {
      if (this.combobox) {
        return 'v-combobox';
      }

      return this.autocomplete ? 'v-autocomplete' : 'v-select';
    },

    content() {
      return this.$refs.select.content;
    },

    rules() {
      return {
        required: this.required,
      };
    },
  },
};
</script>

<style lang="scss">
$selectIconWidth: 24px;

.c-select-field {
  &--ellipsis {
    .v-select__selections {
      width: calc(100% - #{$selectIconWidth});
      flex-wrap: nowrap;
    }
  }
}
</style>