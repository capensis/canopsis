<template>
  <c-advanced-data-table
    :headers="headers"
    :items="entities"
    :loading="pending"
    :options="options"
    :total-items="meta.total_count"
    item-key="_id"
    search
    advanced-pagination
    @update:options="updateOptions"
  />
</template>

<script>
import { useI18n } from '@/hooks/i18n';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { usePbehaviorEntities } from '@/hooks/store/modules/pbehavior-entities';

export default {
  props: {
    pbehavior: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { fetchPbehaviorEntitiesListWithoutStore } = usePbehaviorEntities();

    const {
      data: entities,
      meta,
      pending,
      fetchList,
      options,
      updateOptions,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: ({ params }) =>
        fetchPbehaviorEntitiesListWithoutStore({
          id: props.pbehavior._id,
          params,
        }),
    });

    return {
      t,
      entities,
      meta,
      pending,
      fetchList,
      options,
      updateOptions,
    };
  },

  computed: {
    headers() {
      return [
        {
          text: this.t('common.id'),
          sortable: false,
          value: '_id',
        },
        {
          text: this.t('common.name'),
          sortable: false,
          value: 'name',
        },
        {
          text: this.t('common.type'),
          sortable: false,
          value: 'type',
        },
      ];
    },
  },

  watch: {
    pbehavior() {
      return this.fetchList();
    },
  },

  mounted() {
    return this.fetchList();
  },
};
</script>
