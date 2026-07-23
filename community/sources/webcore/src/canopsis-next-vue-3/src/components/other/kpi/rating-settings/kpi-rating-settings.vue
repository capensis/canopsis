<template>
  <c-layout column>
    <span class="pl-4 text-subtitle-1">
      {{ t('kpi.ratingSettings.helpInformation') }}
    </span>
    <kpi-rating-settings-list
      :options="options"
      :rating-settings="ratingSettings"
      :pending="pending"
      :total-items="meta.total_count"
      :updatable="hasUpdateAccess"
      @change-selected="changeSelectedRatingSettings"
      @update:options="updateOptions"
    />
  </c-layout>
</template>

<script>
import { onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { USER_PERMISSIONS } from '@/constants';

import { getQueryForList } from '@/helpers/entities/shared/query';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useRatingSettings } from '@/hooks/store/modules/rating-settings';

import KpiRatingSettingsList from './kpi-rating-settings-list.vue';

export default {
  components: { KpiRatingSettingsList },

  setup() {
    const { t } = useI18n();
    const popups = usePopups();
    const {
      ratingSettings,
      ratingSettingsMeta,
      ratingSettingsPending,
      fetchRatingSettingsList,
      bulkUpdateRatingSettings,
    } = useRatingSettings();
    const { hasUpdateAccess } = useCRUDPermissions(USER_PERMISSIONS.technical.kpiRatingSettings);

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
        fetchRatingSettingsList({
          params: getQueryForList(fetchQuery),
        }),
    });

    const changeSelectedRatingSettings = async (changedRatingSettings) => {
      await bulkUpdateRatingSettings({
        data: changedRatingSettings.map((ratingSetting) => ({
          id: ratingSetting.id,
          enabled: ratingSetting.enabled,
        })),
      });

      popups.success({ text: t('success.default') });
      await fetchList();
    };

    onMounted(fetchList);

    return {
      t,
      ratingSettings,
      meta: ratingSettingsMeta,
      pending: ratingSettingsPending,
      options,
      hasUpdateAccess,
      updateOptions,
      fetchList,
      changeSelectedRatingSettings,
    };
  },
};
</script>
