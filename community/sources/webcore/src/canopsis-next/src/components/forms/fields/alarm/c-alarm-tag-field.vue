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
    :item-text="itemText"
    :item-value="itemValue"
    :multiple="multiple"
    class="c-alarm-tag-field"
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
      <c-chip
        v-if="!showCount || index < showCount"
        :color="item.color"
        :title="item[itemText]"
        class="c-alarm-tag-field__tag px-2"
        closable
        ellipsis
        @close="removeItemFromSelectedItemsByIndex(index)"
      >
        {{ item[itemText] }}
      </c-chip>
      <span v-else-if="index === showCount">+{{ selectedItems.length - showCount }} {{ $t('common.more') }}</span>
      <span v-else />
    </template>
    <template #item="{ item, attrs, on, parent }">
      <v-list-item
        class="c-alarm-tag-field__list-item"
        v-bind="attrs"
        v-on="on"
      >
        <v-list-item-action v-if="multiple">
          <v-checkbox
            :input-value="attrs.inputValue"
            :color="parent.color"
          />
        </v-list-item-action>
        <v-list-item-content class="c-word-break-all">
          {{ item[itemText] }}
        </v-list-item-content>
      </v-list-item>
    </template>
    <template v-if="$slots['no-data']" #no-data="">
      <slot name="no-data" />
    </template>
  </c-lazy-search-field>
</template>

<script>
import { computed, toRef, watch } from 'vue';

import { PAGINATION_LIMIT } from '@/config';

import { useAlarmTag } from '@/hooks/store/modules/alarm-tag';
import { useLazySearch } from '@/hooks/form/lazy-search';
import { useAlarmTagLabel } from '@/hooks/store/modules/alarm-tag-label';

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
    multiple: {
      type: Boolean,
      required: false,
    },
    onlyLabels: {
      type: Boolean,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { fetchAlarmTagsListWithoutStore } = useAlarmTag();
    const { fetchAlarmTagsLabelsListWithoutStore } = useAlarmTagLabel();

    const itemText = computed(() => (props.onlyLabels ? '_id' : 'value'));
    const itemValue = computed(() => (props.onlyLabels ? '_id' : 'value'));
    const idParamsKey = computed(() => (props.onlyLabels ? 'ids' : 'values'));
    const fetchHandler = computed(() => (
      props.onlyLabels ? fetchAlarmTagsLabelsListWithoutStore : fetchAlarmTagsListWithoutStore
    ));

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
      idParamsKey,
      fetchHandler,
      idKey: itemValue,
      value: toRef(props, 'value'),
      addable: toRef(props, 'addable'),
      multiple: true,
    }, emit);

    watch(() => props.onlyLabels, fetchItems);

    return {
      selectedItems,
      items,
      wholePending,
      hasMoreItems,
      itemText,
      itemValue,
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
