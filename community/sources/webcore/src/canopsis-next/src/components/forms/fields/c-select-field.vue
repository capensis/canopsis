<template>
  <component
    v-bind="$attrs"
    :is="component"
    v-validate="rules"
    ref="select"
    :value="value"
    :items="items"
    :class="{ 'c-select-field--ellipsis': ellipsis }"
    :item-text="itemText"
    :item-value="itemValue"
    :name="name"
    :error-messages="errors.collect(name)"
    :menu-props="computedMenuProps"
    class="c-select-field"
    v-on="$listeners"
  >
    <template
      v-if="$scopedSlots.selection || ellipsis"
      #selection="props"
    >
      <slot
        name="selection"
        v-bind="props"
      >
        <span class="text-truncate">
          {{ getItemText(props.item) }}
        </span>
      </slot>
    </template>
    <template
      v-if="$scopedSlots.item || hasHeaders"
      #item="{ item, attrs, on }"
    >
      <slot
        :item="item"
        :attrs="attrs"
        :on="on"
        name="item"
      >
        <v-subheader
          v-if="item.header"
          :key="item.header"
          class="c-select-field__header"
        >
          {{ item.header }}
        </v-subheader>
        <v-list-item
          v-else
          v-bind="attrs"
          :class="{ 'c-select-field__item': hasHeaders }"
          v-on="on"
        >
          <v-list-item-content>
            <v-list-item-title>{{ getItemText(item) }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </slot>
    </template>
    <template
      v-if="$scopedSlots['append-item']"
      #append-item=""
    >
      <slot name="append-item" />
    </template>
    <template v-if="$slots['no-data']" #no-data="">
      <slot name="no-data" />
    </template>
  </component>
</template>

<script>
import { computed, ref } from 'vue';
import { Validator } from 'vee-validate';
import { isFunction, isObject } from 'lodash';

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
    items: {
      type: Array,
      default: () => [],
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
      type: [String, Function],
      default: 'value',
    },
    name: {
      type: String,
      default: 'value',
    },
    menuProps: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const select = ref(null);

    const component = computed(() => {
      if (props.combobox) {
        return 'v-combobox';
      }

      return props.autocomplete ? 'v-autocomplete' : 'v-select';
    });

    const rules = computed(() => ({
      required: props.required,
    }));

    const hasHeaders = computed(() => props.items.some(item => !!item?.header));

    const computedMenuProps = computed(() => {
      if (!hasHeaders.value) {
        return props.menuProps;
      }

      const contentClass = [
        props.menuProps.contentClass,
        'c-select-field-menu--with-headers',
      ].filter(Boolean).join(' ');

      return {
        ...props.menuProps,
        contentClass,
      };
    });

    const getItemText = (item) => {
      if (isFunction(props.itemText)) {
        return props.itemText(item);
      }

      return isObject(item) ? item[props.itemText] : item;
    };

    return {
      select,

      component,
      rules,
      hasHeaders,
      computedMenuProps,

      getItemText,
    };
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

.c-select-field-menu--with-headers {
  .c-select-field__header {
    font-weight: 700;
    height: auto;
    min-height: 32px;
  }

  .c-select-field__item {
    padding-left: 32px !important;
  }
}
</style>
