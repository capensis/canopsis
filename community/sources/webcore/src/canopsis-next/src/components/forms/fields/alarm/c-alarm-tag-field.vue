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
    item-text="value"
    item-value="value"
    multiple
    chips
    dense
    clearable
    return-object
    @input="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  >
    <template #selection="{ item, index }">
      <c-alarm-action-chip
        v-if="!showCount || index < showCount"
        :color="item.color"
        :title="item.value"
        class="c-alarm-tag-field__tag px-2"
        closable
        ellipsis
        @close="removeItemFromSelectedItemsByIndex(index)"
      >
        {{ item.value }}
      </c-alarm-action-chip>
      <span v-else-if="index === showCount">+{{ selectedItems.length - showCount }} {{ $t('common.more') }}</span>
      <span v-else />
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
          {{ item.value }}
        </v-list-item-content>
      </v-list-item>
    </template>
  </c-lazy-search-field>
</template>

<script>
import { toRef } from 'vue';

import { PAGINATION_LIMIT } from '@/config';

import { useAlarmTag } from '@/hooks/store/modules/alarm-tag';
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
    addable: {
      type: Boolean,
      default: false,
    },
    showCount: {
      type: Number,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { fetchAlarmTagsListWithoutStore } = useAlarmTag();

    const {
      selectedItems,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      removeItemFromSelectedItemsByIndex,
      updateSearch,
    } = useLazySearch({
      value: toRef(props, 'value'),
      addable: toRef(props, 'addable'),
      idKey: 'value',
      idParamsKey: 'values',
      fetchHandler: fetchAlarmTagsListWithoutStore,
    }, emit);

    return {
      selectedItems,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      updateSearch,
      removeItemFromSelectedItemsByIndex,
    };
  },
};
</script>

<style lang="scss">
$selectIconsWidth: 56px;

.c-alarm-tag-field {
  .v-select__selections {
    width: calc(100% - 56px);
  }

  &__tag {
    max-width: 100%;
  }

  &__list {
    max-width: 400px;
  }

  &__list-item .v-list-item {
    height: unset !important;
  }
}
</style>
