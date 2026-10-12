<template>
  <c-lazy-search-field
    :value="selectedItem"
    :label="$t('common.category')"
    :items="items"
    :loading="wholePending || creating"
    :readonly="creating"
    :has-more="hasMoreItems"
    :name="name"
    :hide-details="hideDetails"
    :clearable="isClearable"
    :required="required"
    class="mt-0"
    item-text="name"
    item-value="_id"
    autocomplete
    return-object
    @update:value="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  >
    <template v-if="addable" #prepend-item>
      <div
        class="c-entity-category-field__create-wrapper"
        @mousedown.stop
        @click.stop
        @keydown.stop
      >
        <c-text-field
          :model-value="newCategory"
          :label="$t('service.createCategory')"
          :placeholder="$t('service.createCategory')"
          :name="newCategoryFieldName"
          :unique="isCategoryNameUnique"
          :disabled="creating"
          class="c-entity-category-field__create px-3 pt-1 pb-3"
          hide-details
          hide-label
          @update:model-value="newCategory = $event"
          @keydown.enter.stop.prevent="createCategory"
          @blur="clearCategory"
        >
          <template #append-inner>
            <c-help-icon :text="$t('service.createCategoryHelp')" icon="help" left />
          </template>
        </c-text-field>
      </div>
    </template>
  </c-lazy-search-field>
</template>

<script>
import { computed, ref, toRef } from 'vue';

import { useFormContext } from '@/hooks/form';
import { useLazySearch } from '@/hooks/form/lazy-search';
import { useEntityCategory } from '@/hooks/store/modules/entity-category';

export default {
  props: {
    category: {
      type: [Object, String],
      default: '',
    },
    name: {
      type: String,
      default: 'category',
    },
    addable: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: undefined,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['input', 'update:category', 'update:selected-items'],

  setup(props, { emit }) {
    const form = useFormContext();
    const { createEntityCategory, fetchCategoriesListWithoutStore } = useEntityCategory();
    const newCategory = ref('');
    const creating = ref(false);
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
        value: toRef(props, 'category'),
        idKey: '_id',
        idParamsKey: 'ids',
        fetchHandler: fetchCategoriesListWithoutStore,
        returnObject: true,
        attachValue: true,
      },
      emit,
    );

    const newCategoryFieldName = computed(() => `${props.name}.create`);
    const isClearable = computed(() => props.clearable ?? !props.required);
    const categoryNames = computed(() =>
      items.value
        .map(({ name: categoryName = '' }) => categoryName.trim().toLowerCase())
        .filter(Boolean),
    );
    const isCategoryNameUnique = value =>
      !categoryNames.value.includes(String(value).trim().toLowerCase());

    const clearCategory = () => {
      if (creating.value) {
        return;
      }

      newCategory.value = '';
      form.resetField(newCategoryFieldName.value, { value: '' });
    };

    const createCategory = async () => {
      const name = newCategory.value.trim();

      if (!name || creating.value) {
        return;
      }

      const { valid } = await form.validateField(newCategoryFieldName.value);

      if (!valid) {
        return;
      }

      creating.value = true;

      try {
        const category = await createEntityCategory({ data: { name } });

        await fetchItems();

        changeSelectedItems(category);
        newCategory.value = '';
        form.resetField(newCategoryFieldName.value, { value: '' });
      } finally {
        creating.value = false;
      }
    };

    return {
      newCategory,
      creating,
      newCategoryFieldName,
      isClearable,
      selectedItem,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      updateSearch,
      isCategoryNameUnique,
      createCategory,
      clearCategory,
    };
  },
};
</script>
