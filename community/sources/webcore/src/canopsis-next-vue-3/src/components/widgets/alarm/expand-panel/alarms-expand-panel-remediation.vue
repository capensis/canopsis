<template>
  <remediation-instruction-executions-list
    :executions="executions"
    :options="options"
    :total-items="meta.total_count"
    :pending="pending"
    @update:options="updateOptions"
  />
</template>

<script>
import { onMounted } from 'vue';

import { useChildObserver } from '@/hooks/observer';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useAlarm } from '@/hooks/store/modules/alarm';

import RemediationInstructionExecutionsList from '@/components/other/remediation/instruction-execute/remediation-instruction-executions-list.vue';

export default {
  components: { RemediationInstructionExecutionsList },

  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { fetchExecutionsWithoutStore } = useAlarm();

    const {
      data: executions,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: ({ params }) =>
        fetchExecutionsWithoutStore({
          id: props.alarm._id,
          params,
        }),
    });

    const fetchListWithErrorHandling = async () => {
      try {
        await fetchList();
      } catch (err) {
        console.error(err);
      }
    };

    useChildObserver({ key: '$periodicRefresh', handler: fetchListWithErrorHandling });
    onMounted(fetchListWithErrorHandling);

    return {
      executions,
      meta,
      options,
      pending,
      updateOptions,
    };
  },
};
</script>
