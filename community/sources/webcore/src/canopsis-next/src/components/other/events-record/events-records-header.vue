<template>
  <div class="events-records-header pa-3 ga-3 d-flex flex-column align-center">
    <h3>{{ t('eventsRecord.subheader') }}</h3>

    <c-layout
      v-for="message in messages"
      :key="message"
      gap="4"
      align="center"
    >
      <span class="font-italic">
        {{ message }}
      </span>
      <v-progress-circular color="primary" width="3" indeterminate />
    </c-layout>

    <div>
      <v-btn :disabled="isLaunchDisabled" color="warning" @click="$emit('start:recording')">
        {{ t('eventsRecord.launch') }}
      </v-btn>
    </div>

    <span class="text-caption text-medium-emphasis">
      <strong>{{ t('common.limit') }}:</strong>&nbsp;
      {{ t('eventsRecord.concurrentRecordings', { limit }) }}
    </span>
  </div>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    recordings: {
      type: Array,
      default: () => [],
    },
    resendings: {
      type: Array,
      default: () => [],
    },
    limit: {
      type: Number,
      default: 0,
    },
  },

  emits: ['start:recording'],

  setup(props) {
    const { t } = useI18n();
    const messages = computed(() => {
      const result = [];

      if (props.recordings?.length) {
        result.push(t('eventsRecord.inProgress'));
      }

      if (props.resendings?.length) {
        result.push(t('eventsRecord.resendingInProgress'));
      }

      return result;
    });

    const isLaunchDisabled = computed(
      () => !props.limit || props.recordings?.length >= props.limit,
    );

    return {
      t,
      messages,
      isLaunchDisabled,
    };
  },
};
</script>

<style lang="scss" scoped>
.events-records-header {
  width: 100%;
  text-align: center;
}
</style>
