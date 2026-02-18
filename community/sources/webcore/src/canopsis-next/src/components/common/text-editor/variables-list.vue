<template>
  <c-list
    v-field="value"
    v-bind="$attrs"
    :items="items"
    :item-value="itemValue"
    :item-text="itemText"
    :show-value="showValue"
    :hide-empty-value="hideEmptyValue"
    v-on="$listeners"
  >
    <template v-if="$scopedSlots.prepend" #prepend>
      <slot :items="items" name="prepend" />
    </template>
    <template #item-title="{ item }">
      <v-layout class="gap-4" justify-space-between>
        <v-layout v-if="item[itemText] || item.alias" align-center>
          <v-list-item-mask v-if="item[itemText]" :text="item[itemText]" :mask="search" />
          <v-icon
            v-if="item.alias"
            class="ml-1"
            color="primary"
            small
          >
            alternate_email
          </v-icon>
        </v-layout>
        <span
          v-if="showValue && (!hideEmptyValue || !isUndefined(item[itemValue]))"
          class="grey--text lighten-1"
        >
          {{ item[itemValue] }}
        </span>
        <span
          v-if="item.defined"
          class="grey--text lighten-1"
        >
          {{ $t('common.defined') }}
        </span>
      </v-layout>
    </template>
  </c-list>
</template>
<script>
import { isUndefined } from 'lodash';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Array, String, Number],
      default: () => [],
    },
    search: {
      type: String,
      default: '',
    },
    items: {
      type: Array,
      default: () => [],
    },
    itemValue: {
      type: String,
      default: 'value',
    },
    itemText: {
      type: String,
      default: 'text',
    },
    showValue: {
      type: Boolean,
      default: false,
    },
    hideEmptyValue: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    return {
      isUndefined,
    };
  },
};
</script>
