<template>
  <c-advanced-data-table
    :items="comments"
    :headers="headers"
    :loading="pending"
    :options="options"
    :total-items="totalItems"
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #created="{ item }">
      {{ dateToString(item.created) }}
    </template>
    <template #rating="{ item }">
      <rating-field :model-value="item.rating" readonly />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed, onMounted, ref, watch } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useRemediationInstructionStats } from '@/hooks/store/modules/remediation-instruction-stats';

import RatingField from '@/components/forms/fields/rating-field.vue';

export default {
  components: { RatingField },

  props: {
    remediationInstruction: {
      type: Object,
      required: true,
    },
    interval: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { fetchRemediationInstructionStatsCommentsWithoutStore } =
      useRemediationInstructionStats();
    const totalItems = ref(0);

    const {
      data: comments,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: ({ params }) => fetchRemediationInstructionStatsCommentsWithoutStore({
          id: props.remediationInstruction._id,
          params: {
            ...params,
            from: props.interval.from,
            to: props.interval.to,
          },
        }),
    });

    watch(meta, (value) => {
      totalItems.value = value?.total_count ?? 0;
    });
    watch(() => props.interval, fetchList);

    const headers = computed(() => [
      { text: t('common.date'), value: 'created', sortable: false },
      { text: t('common.username'), value: 'user.display_name', sortable: false },
      { text: t('common.rating'), value: 'rating', sortable: false },
      { text: t('common.comment', 2), value: 'comment', sortable: false },
    ]);
    const dateToString = value => convertDateToString(value);

    onMounted(fetchList);

    return {
      comments,
      pending,
      options,
      totalItems,
      headers,
      updateOptions,
      dateToString,
    };
  },
};
</script>
