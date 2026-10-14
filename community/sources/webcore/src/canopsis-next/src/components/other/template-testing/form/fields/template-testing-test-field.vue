<template>
  <c-lazy-search-field
    :value="selectedItems?.[0]"
    :items="items"
    :label="t('templateTesting.testName')"
    :loading="wholePending"
    :name="name"
    :has-more="hasMoreItems"
    :required="required"
    :return-object="returnObject"
    item-text="name"
    item-value="_id"
    autocomplete
    @update:value="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  />
</template>

<script>
import { merge } from 'lodash';
import { toRef } from 'vue';

import { useLazySearch } from '@/hooks/form/lazy-search';
import { useI18n } from '@/hooks/i18n';
import { useTemplateTesting } from '@/hooks/store/modules/template-testing';

export default {
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'name',
    },
    ruleId: {
      type: String,
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:value', 'update:selected-items'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { fetchTemplateTestingTestsListWithoutStore } = useTemplateTesting();

    const fetchHandler = ({ params }) => fetchTemplateTestingTestsListWithoutStore({
        params: merge(params, { rule: props.ruleId }),
      });

    const {
      selectedItems,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      updateSearch,
    } = useLazySearch(
      {
        fetchHandler,
        value: toRef(props, 'value'),
        returnObject: toRef(props, 'returnObject'),
        idParamsKey: 'ids',
        idKey: '_id',
        attachValue: true,
      },
      emit,
    );

    return {
      t,
      selectedItems,
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
