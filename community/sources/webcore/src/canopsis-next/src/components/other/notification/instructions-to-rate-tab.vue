<template>
  <remediation-instruction-stats-list
    :remediation-instruction-stats="items"
    :pending="pending"
    :options="options"
    :total-items="meta.total_count"
    :accumulated-before="meta.accumulated_before"
    :interval="interval"
    @rate="showRateInstructionModal"
    @update:options="$emit('update:options', $event)"
  />
</template>

<script>
import { toRef } from 'vue';

import {
  useRemediationInstructionStatsRate,
} from '@/components/other/remediation/instruction-stats/hooks/remediation-instruction-stats';

import RemediationInstructionStatsList from '@/components/other/remediation/instruction-stats/remediation-instruction-stats-list.vue';

import { useNotificationActiveId } from './hooks/notifications';

export default {
  components: {
    RemediationInstructionStatsList,
  },
  props: {
    items: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    meta: {
      type: Object,
      default: () => {},
    },
    options: {
      type: Object,
      default: () => {},
    },
    interval: {
      type: Object,
      default: () => {},
    },
    activeId: {
      type: String,
      required: false,
    },
  },
  setup(props, { emit }) {
    const refresh = () => emit('refresh');

    const { showRateInstructionModal } = useRemediationInstructionStatsRate(refresh);

    useNotificationActiveId({
      activeId: toRef(props, 'activeId'),
      items: toRef(props, 'items'),
      action: showRateInstructionModal,
    });

    return {
      showRateInstructionModal,
    };
  },
};
</script>
