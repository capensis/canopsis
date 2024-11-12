<template>
  <c-lazy-search-field
    :value="selectedItems"
    :items="items"
    :label="label || $tc('common.tag')"
    :loading="wholePending"
    :disabled="disabled"
    :name="name"
    :menu-props="{ contentClass: 'c-alarm-tag-field__list' }"
    :has-more="hasMoreItems"
    :required="required"
    :autocomplete="!combobox"
    :hide-details="!required"
    :hide-selected="hideSelected"
    class="c-alarm-tag-field"
    item-text="_id"
    item-value="_id"
    multiple
    chips
    dense
    clearable
    return-object
    @input="changeSelectedItems"
    @fetch="fetchTags"
    @fetch:more="fetchMoreTags"
    @update:search="updateSearch"
  >
    <template #selection="{ item, index }">
      <c-alarm-action-chip
        :color="item.color"
        :title="item._id"
        :text-color="item.color ? 'white' : 'black'"
        class="c-alarm-tag-field__tag px-2"
        closable
        ellipsis
        @close="removeItemFromSelectedItemsByIndex(index)"
      >
        {{ item._id ?? item }}
      </c-alarm-action-chip>
    </template>
    <template #item="{ item, attrs, on, parent }">
      <v-list-item
        class="c-alarm-tag-field__list-item"
        v-bind="attrs"
        v-on="on"
      >
        <v-list-item-action>
          <v-checkbox
            :input-value="attrs.inputValue"
            :color="parent.color"
          />
        </v-list-item-action>
        <v-list-item-content class="c-word-break-all">
          {{ item._id }}
        </v-list-item-content>
      </v-list-item>
    </template>
    <template v-if="$slots['no-data']" #no-data="">
      <slot name="no-data" />
    </template>
  </c-lazy-search-field>
</template>

<script>
import { toRef } from 'vue';

import { PAGINATION_LIMIT } from '@/config';

import { useAlarmTagLabel } from '@/hooks/store/modules/alarm-tag-label';
import { useLazySearch } from '@/hooks/form/lazy-search';

export default {
  props: {
    value: {
      type: [Array],
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'tag',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    limit: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
    combobox: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    hideSelected: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const { fetchAlarmTagsLabelsListWithoutStore } = useAlarmTagLabel();

    const {
      selectedItems,
      items,
      wholePending,
      hasMoreItems,
      fetchTags,
      fetchMoreTags,
      changeSelectedItems,
      removeItemFromSelectedItemsByIndex,
      updateSearch,
    } = useLazySearch({
      value: toRef(props, 'value'),
      idKey: '_id',
      idParamsKey: 'ids',
      fetchHandler: fetchAlarmTagsLabelsListWithoutStore,
      addable: true,
    }, emit);

    return {
      selectedItems,
      items,
      wholePending,
      hasMoreItems,
      fetchTags,
      fetchMoreTags,
      changeSelectedItems,
      updateSearch,
      removeItemFromSelectedItemsByIndex,
    };
  },
};
</script>
