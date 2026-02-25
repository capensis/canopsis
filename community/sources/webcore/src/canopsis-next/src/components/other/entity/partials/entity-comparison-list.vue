<template>
  <c-collapse-panel class="c-alternative-bg-panel" expanded>
    <template #header>
      <span class="font-weight-medium text-uppercase">{{ title }} - {{ meta.total_count ?? '' }}</span>
    </template>
    <entities-list-table-with-pagination
      :widget="widget"
      :columns="widget.parameters.widgetColumns"
      :entities="entities"
      :meta="meta"
      :loading="pending"
      :query="query"
      selectable
      @update:query="updateQuery"
    />
  </c-collapse-panel>
</template>

<script>
import { computed, ref, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { WIDGET_TYPES, ENTITY_FIELDS } from '@/constants';

import { convertQueryToRequest } from '@/helpers/query';
import { formToWidget, widgetToForm } from '@/helpers/entities/widget/form';

import { useI18n } from '@/hooks/i18n';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useEntity } from '@/hooks/store/modules/entity';

import EntitiesListTableWithPagination from '@/components/widgets/context/partials/entities-list-table-with-pagination.vue';

export default {
  components: {
    EntitiesListTableWithPagination,
  },
  props: {
    title: {
      type: String,
      default: '',
    },
    initialQuery: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const { t } = useI18n();

    const entities = ref([]);
    const meta = ref({});

    const { fetchContextEntitiesListWithoutStore } = useEntity();
    const {
      pending,
      query,
      updateQuery,
      fetchHandlerWithQuery: fetchList,
    } = usePendingWithLocalQuery({
      initialQuery: { page: 1, itemsPerPage: PAGINATION_LIMIT, ...props.initialQuery },
      fetchHandler: async (fetchQuery) => {
        const response = await fetchContextEntitiesListWithoutStore({
          params: {
            ...convertQueryToRequest(fetchQuery),

            entity_pattern: fetchQuery.entity_pattern,
            negative_entity_pattern: fetchQuery.negative_entity_pattern,
          },
        });

        entities.value = response.data;
        meta.value = response.meta;
      },
    });

    const widget = computed(() => formToWidget(widgetToForm({
      type: WIDGET_TYPES.context,
      parameters: {
        widgetColumns: [
          {
            value: ENTITY_FIELDS.name,
            label: t('common.name'),
          },
          {
            value: ENTITY_FIELDS.type,
            label: t('common.type'),
          },
        ],
      },
    })));

    onMounted(() => fetchList());

    return {
      entities,
      meta,
      widget,
      pending,
      query,
      updateQuery,
      fetchList,
    };
  },
};
</script>
