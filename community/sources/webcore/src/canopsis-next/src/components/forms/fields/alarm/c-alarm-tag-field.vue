<template>
  <c-lazy-search-field
    :value="selectedItems"
    :items="items"
    :label="label || $t('common.tag')"
    :loading="wholePending"
    :disabled="disabled"
    :name="name"
    :menu-props="{ contentClass: 'c-alarm-tag-field__list' }"
    :has-more="hasMoreItems"
    :required="required"
    :autocomplete="!combobox"
    :hide-details="!required"
    :hide-selected="hideSelected"
    :hide-no-data="false"
    :item-text="itemText"
    :item-value="itemValue"
    :multiple="multiple"
    :show-count="showCount"
    class="c-alarm-tag-field"
    chips
    clearable
    return-object
    @update:value="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  >
    <template #selection="{ item, index }">
      <c-chip
        :color="getSelectionColor(item)"
        :text-color="getSelectionTextColor(item)"
        :title="getSelectionText(item)"
        class="c-alarm-tag-field__tag px-2"
        closable
        @close="removeItemFromSelectedItemsByIndex(index)"
      >
        {{ getSelectionText(item) }}
      </c-chip>
    </template>
    <template #item="{ item, props: itemProps, search }">
      <v-list-item class="c-alarm-tag-field__list-item text-break" v-bind="itemProps">
        <template v-if="multiple" #prepend="{ isSelected }">
          <v-checkbox-btn :model-value="isSelected" color="primary" />
        </template>

        <v-list-item-title class="c-word-break-all">
          <v-list-item-mask :text="String(getSelectionText(item) ?? '')" :mask="search" />
        </v-list-item-title>
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

import { getMostReadableTextColor } from '@/helpers/color';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useAlarmTag } from '@/hooks/store/modules/alarm-tag';
import { useAlarmTagLabel } from '@/hooks/store/modules/alarm-tag-label';

export default {
  emits: ['update:value', 'update:selected-items'],
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
    const fetchHandler = computed(
      () => (props.onlyLabels ? fetchAlarmTagsLabelsListWithoutStore : fetchAlarmTagsListWithoutStore),
    );

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
    } = useLazySearch(
      {
        idParamsKey,
        fetchHandler,
        idKey: itemValue,
        value: toRef(props, 'value'),
        addable: toRef(props, 'addable'),
        multiple: true,
        attachValue: true,
      },
      emit,
    );

    const getSelectionText = item => item?.[itemText.value] ?? item;
    const getSelectionColor = item => item?.color;
    const getSelectionTextColor = item =>
      getMostReadableTextColor(getSelectionColor(item), { level: 'AA', size: 'large' });

    watch(() => props.onlyLabels, () => fetchItems());

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
      getSelectionText,
      getSelectionColor,
      getSelectionTextColor,
      removeItemFromSelectedItemsByIndex,
    };
  },
};
</script>

<style lang="scss">
.c-alarm-tag-field {
  .v-select__selection,
  .v-autocomplete__selection {
    width: auto;
    max-width: calc(100% - 56px);
    height: 24px;
    margin: 0;
  }

  &__tag.c-chip.v-chip {
    display: inline-flex;
    align-items: center;
    max-width: 100%;
    height: 24px !important;
    min-height: 24px !important;
    margin: 0;

    .c-chip__text {
      display: flex;
      align-items: center;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
      line-height: 1;
    }

    .v-chip__content {
      display: flex !important;
      align-items: center;
      height: 100% !important;
    }
  }

  &__list {
    max-width: 400px;
  }

  &__list-item.v-list-item {
    min-height: unset !important;
  }
}
</style>
