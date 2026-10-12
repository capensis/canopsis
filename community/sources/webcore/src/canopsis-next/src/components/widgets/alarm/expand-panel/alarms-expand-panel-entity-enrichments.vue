<template>
  <entity-infos-logs-list
    :items="items"
    :options="options"
    :total-items="meta.total_count"
    :pending="pending"
    @update:options="updateOptions"
  />
</template>

<script>
import { onMounted } from 'vue';

import { useEntity } from '@/hooks/store/modules/entity';
import { usePeriodicRefreshChild } from '@/hooks/periodic-refresh';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';

import EntityInfosLogsList from '@/components/other/entity/entity-infos-logs-list.vue';

export default {
  components: {
    EntityInfosLogsList,
  },
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { fetchEntityInfosLogsListWithoutStore } = useEntity();

    const {
      data: items,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: rest => fetchEntityInfosLogsListWithoutStore({
        ...rest,

        params: {
          ...rest.params,

          _id: props.alarm?.entity?._id,
        },
      }),
    });

    usePeriodicRefreshChild({ handler: fetchList });

    onMounted(fetchList);

    return {
      items,
      meta,
      pending,
      options,
      updateOptions,
    };
  },
};
</script>
