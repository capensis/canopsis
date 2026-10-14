<template>
  <c-list
    v-bind="$attrs"
    :value="value"
    :items="items"
    :children-key="childrenKey"
    :item-value="itemValue"
    :item-text="itemText"
    :show-value="showValue"
    :hide-empty-value="hideEmptyValue"
    :dense="dense"
    :return-object="returnObject"
    :clickable-parent="clickableParent"
    @update:value="$emit('update:value', $event)"
    @fetch:more="$emit('fetch:more')"
  >
    <template v-if="$slots['no-data']" #no-data>
      <slot name="no-data" />
    </template>

    <template v-if="$slots.prepend" #prepend="{ items: prependItems }">
      <slot :items="prependItems" name="prepend" />
    </template>

    <template #item-title="{ item }">
      <c-layout gap="4" justify="space-between">
        <c-layout v-if="item[itemText] || item.alias" align="center">
          <v-list-item-mask v-if="item[itemText]" :text="item[itemText]" :mask="search" />

          <c-icon
            v-if="item.alias"
            class="ml-1"
            color="primary"
            size="small"
          >
            alternate_email
          </c-icon>
        </c-layout>

        <span
          v-if="showValue && (!hideEmptyValue || item[itemValue] !== undefined)"
          class="text-medium-emphasis"
        >
          {{ item[itemValue] }}
        </span>

        <span v-if="item.defined" class="text-medium-emphasis">
          {{ $t('common.defined') }}
        </span>
      </c-layout>
    </template>
  </c-list>
</template>

<script>
export default {
  inheritAttrs: false,

  props: {
    value: {
      type: [Array, String, Number, Object, Boolean],
      default: () => [],
    },
    search: {
      type: [String, Number],
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
    childrenKey: {
      type: String,
      default: 'items',
    },
    showValue: {
      type: Boolean,
      default: false,
    },
    hideEmptyValue: {
      type: Boolean,
      default: false,
    },
    dense: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    clickableParent: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:value', 'fetch:more'],
};
</script>
