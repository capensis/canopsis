<template>
  <c-lazy-search-field
    :value="selectedItems?.[0]"
    :items="items"
    :label="label || t('templateTesting.testName')"
    :loading="wholePending"
    :name="name"
    :has-more="hasMoreItems"
    :required="required"
    :clearable="clearable"
    :disabled="disabled"
    :menu-props="menuProps"
    item-text="name"
    item-value="_id"
    return-object
    autocomplete
    @update:value="changeSelectedItems"
    @fetch="fetchItems"
    @fetch:more="fetchMoreItems"
    @update:search="updateSearch"
  >
    <template v-if="hasReadAccessForTemplateData && hasCreateAccessForTemplateData" #no-data>
      <c-layout justify="center">
        <v-btn color="primary" variant="outlined" @click="showCreateTemplateTestingDataModal">
          <v-icon class="mr-2"> add_circle </v-icon>

          {{ t('templateTesting.addData') }}
        </v-btn>
      </c-layout>
    </template>
  </c-lazy-search-field>
</template>

<script>
import { merge } from 'lodash';
import { toRef, watch } from 'vue';

import { USER_PERMISSIONS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useCRUDPermissions } from '@/hooks/auth';
import { useLazySearch } from '@/hooks/form/lazy-search';
import { useTemplateTesting } from '@/hooks/store/modules/template-testing';

import { useTemplateDataModals } from '../../hooks/template-testing-data';

export default {
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
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:value'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const menuProps = { closeOnContentClick: true };

    const {
      hasCreateAccess: hasCreateAccessForTemplateData,
      hasReadAccess: hasReadAccessForTemplateData,
    } = useCRUDPermissions(USER_PERMISSIONS.technical.templateData);

    const { fetchTemplateTestingDataListWithoutStore } = useTemplateTesting();

    const fetchHandler = async ({ params }) => {
      if (!hasReadAccessForTemplateData.value) {
        return { data: [], meta: { page_count: 0 } };
      }

      return fetchTemplateTestingDataListWithoutStore({
        params: merge(params, props.params, { type: props.type }),
      });
    };

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
        idParamsKey: 'ids',
        idKey: '_id',
        attachValue: true,
      },
      emit,
    );

    const { showCreateTemplateTestingDataModal } = useTemplateDataModals({
      refresh: fetchItems,
      type: toRef(props, 'type'),
    });

    watch(() => props.params, fetchItems);

    return {
      t,
      menuProps,
      hasCreateAccessForTemplateData,
      hasReadAccessForTemplateData,
      selectedItems,
      items,
      wholePending,
      hasMoreItems,
      fetchItems,
      fetchMoreItems,
      changeSelectedItems,
      updateSearch,
      showCreateTemplateTestingDataModal,
    };
  },
};
</script>
