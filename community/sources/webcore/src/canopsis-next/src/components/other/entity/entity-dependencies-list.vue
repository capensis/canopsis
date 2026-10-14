<template>
  <entities-list-table-with-pagination
    :widget="widget"
    :entities="entities"
    :pending="pending"
    :meta="meta"
    :query="query"
    :columns="columns"
    selectable
    @update:query="updateQuery"
  >
    <template #toolbar>
      <c-layout class="gap-4 py-4 pt-0" align="end" wrap>
        <div class="flex-grow-1">
          <c-advanced-search
            :attributes="advancedSearchAttributes"
            @submit="updateSearch"
            @reset="resetSearch"
          />
        </div>
        <div v-if="hasAccessToCategory">
          <c-entity-category-field
            :category="query.category"
            hide-details
            @update:category="updateCategory"
          />
        </div>
      </c-layout>
    </template>
  </entities-list-table-with-pagination>
</template>

<script>
import { omit } from 'lodash';
import { computed, onMounted, ref } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { ADVANCED_SEARCH_FIELDS_TO_COMPARISON, USER_PERMISSIONS } from '@/constants';

import { getQueryForList } from '@/helpers/entities/shared/query';
import { prepareQueryWithAdvancedSearch } from '@/helpers/search/advanced-search';

import { useCurrentUserPermissions } from '@/hooks/auth';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useService } from '@/hooks/store/modules/service';

import { useEntityDependenciesAdvancedSearchAttributes } from '@/components/common/search/hooks/advanced-search';

import EntitiesListTableWithPagination from '@/components/widgets/context/partials/entities-list-table-with-pagination.vue';

export default {
  components: { EntitiesListTableWithPagination },

  props: {
    entityId: {
      type: String,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    impact: {
      type: Boolean,
      default: false,
    },
    columns: {
      type: Array,
      default: () => [],
    },
  },

  setup(props) {
    const entities = ref([]);
    const meta = ref({});

    const { checkAccess } = useCurrentUserPermissions();
    const { fetchServiceDependenciesWithoutStore, fetchServiceImpactsWithoutStore } = useService();
    const { attributes: advancedSearchAttributes } =
      useEntityDependenciesAdvancedSearchAttributes();

    const {
      pending,
      query,
      updateQuery,
      fetchHandlerWithQuery: fetchList,
    } = usePendingWithLocalQuery({
      initialQuery: {
        page: 1,
        category: '',
        itemsPerPage: PAGINATION_LIMIT,
        search: '',
        sortBy: [],
        sortDesc: [],
      },
      fetchHandler: async (fetchQuery) => {
        const method = props.impact
          ? fetchServiceImpactsWithoutStore
          : fetchServiceDependenciesWithoutStore;
        const params = {
          ...getQueryForList(fetchQuery),
          with_flags: true,
        };
        const response = await method({
          id: props.entityId,
          params,
        });

        entities.value = response.data;
        meta.value = response.meta;
      },
    });

    const hasAccessToCategory = computed(() =>
      checkAccess(USER_PERMISSIONS.business.context.actions.category),
    );

    const resetSearch = () =>
      updateQuery({
        ...omit(query.value, [...ADVANCED_SEARCH_FIELDS_TO_COMPARISON]),
        page: 1,
      });

    const updateSearch = (search = {}) =>
      updateQuery(prepareQueryWithAdvancedSearch(query.value, search));

    const updateCategory = category =>
      updateQuery({
        ...query.value,
        page: 1,
        category: category?._id,
      });

    onMounted(fetchList);

    return {
      advancedSearchAttributes,
      pending,
      entities,
      meta,
      query,
      hasAccessToCategory,
      resetSearch,
      updateQuery,
      updateSearch,
      updateCategory,
    };
  },
};
</script>
