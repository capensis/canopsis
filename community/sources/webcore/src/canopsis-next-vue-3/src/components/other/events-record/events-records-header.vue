<template>
  <c-layout class="gap-4 mb-4" justify="space-between" align="center">
    <div class="text-medium-emphasis">
      <div>{{ t('eventsRecord.subheader') }}</div>
      <div>{{ t('eventsRecord.concurrentRecordings', { limit }) }}</div>
    </div>

    <c-layout class="gap-4" align="center">
      <v-expand-transition>
        <c-layout v-if="recordings.length" class="gap-2" align="center">
          <span class="font-italic">{{ t('eventsRecord.inProgress') }}</span>
          <v-progress-circular color="primary" width="3" indeterminate />
        </c-layout>
      </v-expand-transition>

      <v-expand-transition>
        <c-layout v-if="resendings.length" class="gap-2" align="center">
          <span class="font-italic">{{ t('eventsRecord.resendingInProgress') }}</span>
          <v-progress-circular color="primary" width="3" indeterminate />
        </c-layout>
      </v-expand-transition>

      <v-btn
        :disabled="!limit || recordings.length >= limit"
        color="primary"
        @click="$emit('start:recording')"
      >
        {{ t('eventsRecord.launch') }}
      </v-btn>
    </c-layout>
  </c-layout>
</template>

<script>
import { useI18n } from '@/hooks/i18n'

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

  setup() {
    const { t } = useI18n()

    return { t }
  },
}
</script>
