<template>
  <c-layout column>
    <span class="pl-4">
      {{ t('kpi.filters.helpInformation') }}
    </span>
    <kpi-filters-list
      :options="options"
      :filters="filters"
      :pending="pending"
      :total-items="meta.total_count"
      :duplicable="hasCreateAccess"
      :removable="hasDeleteAccess"
      :updatable="hasUpdateAccess"
      @edit="showEditFilterModal"
      @duplicate="showDuplicateFilterModal"
      @remove="showDeleteFilterModal"
      @update:options="updateOptions"
    />
  </c-layout>
</template>

<script>
import { onMounted } from 'vue';
import { omit } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { MODALS, USER_PERMISSIONS } from '@/constants';

import { getQueryForList } from '@/helpers/entities/shared/query';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useFilter } from '@/hooks/store/modules/filter';

import KpiFiltersList from './kpi-filters-list.vue';

export default {
  components: { KpiFiltersList },

  setup() {
    const { t } = useI18n();
    const modals = useModals();
    const popups = usePopups();
    const {
      filters,
      filtersMeta,
      filtersPending,
      fetchFiltersList,
      createFilter,
      updateFilter,
      removeFilter,
    } = useFilter();
    const { hasCreateAccess, hasUpdateAccess, hasDeleteAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.kpiFilters,
    );

    const {
      options,
      handler: fetchList,
      updateOptions,
    } = useLocalQueryWithOptions({
      initialQuery: {
        page: 1,
        itemsPerPage: PAGINATION_LIMIT,
        sortBy: [],
        sortDesc: [],
        search: '',
      },
      onUpdate: (fetchQuery) =>
        fetchFiltersList({
          params: getQueryForList(fetchQuery),
        }),
    });

    const showEditFilterModal = (filter) =>
      modals.show({
        name: MODALS.createKpiFilter,
        config: {
          filter,
          title: t('modals.createFilter.edit.title'),
          action: async (data) => {
            const result = await updateFilter({ id: filter._id, data });

            popups.success({ text: t('success.default') });
            await fetchList();

            return result;
          },
        },
      });

    const showDuplicateFilterModal = (filter) =>
      modals.show({
        name: MODALS.createKpiFilter,
        config: {
          filter: omit(filter, ['_id']),
          title: t('modals.createFilter.duplicate.title'),
          action: async (data) => {
            const result = await createFilter({ data });

            popups.success({ text: t('success.default') });
            await fetchList();

            return result;
          },
        },
      });

    const showDeleteFilterModal = (id) =>
      modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            await removeFilter({ id });

            popups.success({ text: t('success.default') });
            await fetchList();
          },
        },
      });

    onMounted(fetchList);

    return {
      t,
      filters,
      meta: filtersMeta,
      pending: filtersPending,
      options,
      hasCreateAccess,
      hasUpdateAccess,
      hasDeleteAccess,
      updateOptions,
      fetchList,
      showEditFilterModal,
      showDuplicateFilterModal,
      showDeleteFilterModal,
    };
  },
};
</script>
