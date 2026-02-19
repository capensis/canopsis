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
    <template #toolbar="">
      <v-flex>
        <c-advanced-search
          :attributes="advancedSearchAttributes"
          @submit="updateSearch"
          @reset="resetSearch"
        />
      </v-flex>
      <v-flex v-if="hasAccessToCategory">
        <c-entity-category-field
          :category="query.category"
          class="mr-3"
          @input="updateCategory"
        />
      </v-flex>
    </template>
  </entities-list-table-with-pagination>
</template>

<script>
import { omit } from 'lodash';
import { ref, computed, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { USER_PERMISSIONS, ADVANCED_SEARCH_FIELDS_TO_COMPARISON } from '@/constants';

import { getQueryForList } from '@/helpers/entities/shared/query';
import { prepareQueryWithAdvancedSearch } from '@/helpers/search/advanced-search';

import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useCurrentUserPermissions } from '@/hooks/auth';
import { useService } from '@/hooks/store/modules/service';

import { useEntityDependenciesAdvancedSearchAttributes } from '@/components/common/search/hooks/advanced-search';

import EntitiesListTableWithPagination from '../../widgets/context/partials/entities-list-table-with-pagination.vue';

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
    const {
      fetchServiceDependenciesWithoutStore,
      fetchServiceImpactsWithoutStore,
    } = useService();
    const { attributes: advancedSearchAttributes } = useEntityDependenciesAdvancedSearchAttributes();

    const {
      pending,
      query,
      updateQuery,
      fetchHandlerWithQuery: fetchList,
    } = usePendingWithLocalQuery({
      initialQuery: {
        page: 1,
        itemsPerPage: PAGINATION_LIMIT,
        search: '',
        sortBy: [],
        sortDesc: [],
      },
      fetchHandler: async (fetchQuery) => {
        const method = props.impact
          ? fetchServiceImpactsWithoutStore
          : fetchServiceDependenciesWithoutStore;

        const params = getQueryForList(fetchQuery);
        params.with_flags = true;

        const { data, meta: responseMeta } = await method({
          id: props.entityId,
          params,
        });

        entities.value = data;
        meta.value = responseMeta;
      },
    });

    const hasAccessToCategory = computed(() => checkAccess(USER_PERMISSIONS.business.context.actions.category));

    /**
     * Resets the search query by removing all search-related parameters and resetting pagination to the first page.
     */
    const resetSearch = () => updateQuery({
      ...omit(query.value, [...ADVANCED_SEARCH_FIELDS_TO_COMPARISON]),

      page: 1,
    });

    /**
     * Updates the search query parameter and resets pagination to the first page.
     *
     * @param {string} search - The search term to filter entities
     */
    const updateSearch = (search = {}) => updateQuery(prepareQueryWithAdvancedSearch(query.value, search));

    /**
     * Updates the category filter in the query and resets pagination to the first page.
     *
     * @param {Object|null} category - The category object with _id property, or null to clear the filter
     */
    const updateCategory = category => updateQuery({
      ...query.value,
      page: 1,
      category: category && category._id,
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
