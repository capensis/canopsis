<template>
  <v-layout>
    <c-action-btn
      v-if="item"
      :tooltip="$t('jobs.actions.editJob')"
      icon="edit"
      @click="$emit('edit', item)"
    />
    <template v-if="itemsForStart.length">
      <c-action-btn
        v-if="shownPlay"
        :tooltip="$t('jobs.actions.startJob')"
        icon="play_arrow"
        @click="$emit('play', item)"
      />
      <c-action-btn
        v-else
        :tooltip="$t('jobs.actions.repeatJob')"
        icon="refresh"
        @click="$emit('retry', item)"
      />
    </template>
    <c-action-btn
      v-if="itemsForPause.length"
      :tooltip="$t('jobs.actions.pauseJob')"
      icon="pause"
      @click="$emit('pause', item)"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { JOB_STATUS } from '@/constants';

export default {
  props: {
    item: {
      type: Object,
      required: false,
    },
    items: {
      type: Array,
      default: () => [],
    },
  },
  setup(props) {
    const itemsForStart = computed(() => props.items.filter(item => item.status === JOB_STATUS.stopped));
    const itemsForPause = computed(() => props.items.filter(item => item.status === JOB_STATUS.stopped));
    const shownPlay = computed(() => itemsForStart.value.length > 0);

    return {
      shownPlay,
      itemsForStart,
      itemsForPause,
    };
  },
};
</script>
