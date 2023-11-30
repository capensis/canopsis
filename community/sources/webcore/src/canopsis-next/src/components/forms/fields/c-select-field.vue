<template>
  <component
    class="c-select-field"
    ref="select"
    v-bind="$attrs"
    v-on="$listeners"
    v-validate="rules"
    :value="value"
    :class="{ 'c-select-field--ellipsis': ellipsis }"
    :is="component"
    :item-text="itemText"
    :item-value="itemValue"
    :name="name"
    :error-messages="errors.collect(name)"
  >
    <template
      v-if="$scopedSlots.selection || ellipsis"
      #selection="props"
    >
      <slot
        name="selection"
        v-bind="props"
      >
        <span class="text-truncate">{{ getItemText(props.item) }}</span>
      </slot>
    </template>
    <template
      v-if="$scopedSlots.item"
      #item="props"
    >
      <slot
        name="item"
        v-bind="props"
      />
    </template>
    <template
      v-if="$scopedSlots['append-item']"
      #append-item=""
    >
      <slot name="append-item" />
    </template>
  </component>
</template>

<script>
import { Validator } from 'vee-validate';
import { isArray, isFunction, isObject } from 'lodash';

export default {
  inject: {
    $validator: {
      default: new Validator(),
    },
  },
  inheritAttrs: false,
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Array, Object, String, Symbol, Number],
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
    name: {
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
  methods: {
    getItemText(item) {
      if (isFunction(this.itemText)) {
        return this.itemText(item);
      }

      return isObject(item) ? item[this.itemText] : item;
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
