<template>
  <c-lazy-search-field
    :value="selectedItem"
    :items="items"
    :label="label || $t('externalData.tableField')"
    :disabled="disabled"
    :required="required"
    :loading="wholePending"
    :has-more="hasMoreItems"
    :name="name"
    item-text="value"
    item-value="value"
    class="c-entity-info-property-key-field"
    return-object
    @update:value="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  />
</template>

<script>
import { toRef } from 'vue';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useService } from '@/hooks/store/modules/service';

export default {
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'name',
    },
  },

  emits: ['update:value', 'update:selected-items'],

  setup(props, { emit }) {
    const { fetchEntityInfosKeysWithoutStore } = useService();

    const {
      selectedItem,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      updateSearch,
    } = useLazySearch(
      {
        value: toRef(props, 'value'),
        idKey: 'value',
        idParamsKey: 'ids',
        fetchHandler: fetchEntityInfosKeysWithoutStore,
        addable: true,
      },
      emit,
    );

    return {
      selectedItem,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      updateSearch,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-entity-info-property-key-field {
  :deep(.v-field__input) {
    min-width: 0;
  }

  :deep(.v-select__selection),
  :deep(.v-select__selection-text) {
    display: block;
    min-width: 0;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}
</style>
