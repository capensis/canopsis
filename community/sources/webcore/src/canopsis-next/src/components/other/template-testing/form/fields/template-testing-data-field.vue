<template>
  <c-lazy-search-field
    :value="selectedItems?.[0]"
    :items="items"
    :label="label ||$tc('templateTesting.testName')"
    :loading="wholePending"
    :name="name"
    :has-more="hasMoreItems"
    :required="required"
    :clearable="clearable"
    :menu-props="menuProps"
    :clear-on-empty-search="false"
    item-text="name"
    item-value="_id"
    return-object
    autocomplete
    @input="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  >
    <template v-if="hasAccessToTemplateTesting" #no-data>
      <v-layout justify-center>
        <v-btn
          color="primary"
          outlined
          @click="showCreateTemplateTestingDataModal"
        >
          <v-icon class="mr-2">
            add_circle
          </v-icon>
          {{ $t('templateTesting.addData') }}
        </v-btn>
      </v-layout>
    </template>
  </c-lazy-search-field>
</template>

<script>
import { merge } from 'lodash';
import { watch, toRef } from 'vue';

import { USER_PERMISSIONS } from '@/constants';

import { useCanPermission } from '@/hooks/auth';
import { useLazySearch } from '@/hooks/form/lazy-search';
import { useTemplateData } from '@/hooks/store/modules/template-data';

import { useTemplateDataModals } from '../../hooks/template-testing-data';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'test',
    },
    type: {
      type: Number,
      required: false,
    },
    params: {
      type: Object,
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const menuProps = { closeOnContentClick: true };

    const hasAccessToTemplateTesting = useCanPermission(USER_PERMISSIONS.technical.templateTesting);

    const { fetchTemplateDataListWithoutStore } = useTemplateData();

    const fetchHandler = ({ params }) => fetchTemplateDataListWithoutStore({
      params: merge(params, props.params, { type: props.type }),
    });

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
      fetchHandler,
      value: toRef(props, 'value'),
      idParamsKey: 'ids',
      idKey: '_id',
      attachValue: true,
    }, emit);

    const { showCreateTemplateTestingDataModal } = useTemplateDataModals({
      refresh: fetchItems,
      type: toRef(props, 'type'),
    });

    watch(() => props.params, fetchItems);

    return {
      menuProps,
      hasAccessToTemplateTesting,
      selectedItems,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      removeItemFromSelectedItemsByIndex,
      updateSearch,
      showCreateTemplateTestingDataModal,
    };
  },
};
</script>
