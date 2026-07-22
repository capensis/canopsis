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
    :item-text="itemText"
    :item-value="itemValue"
    :multiple="multiple"
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
      <v-chip
        v-if="!showCount || index < showCount"
        :color="item.color"
        :title="item[itemText]"
        class="c-alarm-tag-field__tag px-2"
        size="small"
        closable
        @click:close="removeItemFromSelectedItemsByIndex(index)"
      >
        <span class="text-truncate">{{ item[itemText] }}</span>
      </v-chip>
      <span v-else-if="index === showCount">
        +{{ selectedItems.length - showCount }} {{ $t('common.more') }}
      </span>
      <span v-else />
    </template>
    <template #item="{ props: itemProps }">
      <v-list-item class="c-alarm-tag-field__list-item text-break" v-bind="itemProps">
        <template v-if="multiple" #prepend>
          <v-checkbox :model-value="itemProps.active" density="compact" hide-details />
        </template>
      </v-list-item>
    </template>
    <template v-if="$slots['no-data']" #no-data="">
      <slot name="no-data" />
    </template>
  </c-lazy-search-field>
</template>

<script>
import { computed, toRef, watch } from 'vue'

import { PAGINATION_LIMIT } from '@/config'

import { useAlarmTag } from '@/hooks/store/modules/alarm-tag'
import { useLazySearch } from '@/hooks/form/lazy-search'
import { useAlarmTagLabel } from '@/hooks/store/modules/alarm-tag-label'

export default {
  emits: ['update:value'],
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
    const { fetchAlarmTagsListWithoutStore } = useAlarmTag()
    const { fetchAlarmTagsLabelsListWithoutStore } = useAlarmTagLabel()

    const itemText = computed(() => (props.onlyLabels ? '_id' : 'value'))
    const itemValue = computed(() => (props.onlyLabels ? '_id' : 'value'))
    const idParamsKey = computed(() => (props.onlyLabels ? 'ids' : 'values'))
    const fetchHandler = computed(() =>
      props.onlyLabels ? fetchAlarmTagsLabelsListWithoutStore : fetchAlarmTagsListWithoutStore,
    )

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
    )

    watch(() => props.onlyLabels, fetchItems)

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
    }
  },
}
</script>

<style lang="scss">
.c-alarm-tag-field {
  .v-select__selection {
    width: calc(100% - 56px);
  }

  &__tag {
    max-width: 100%;
  }

  &__list {
    max-width: 400px;
  }

  &__list-item.v-list-item {
    min-height: unset !important;
  }
}
</style>
