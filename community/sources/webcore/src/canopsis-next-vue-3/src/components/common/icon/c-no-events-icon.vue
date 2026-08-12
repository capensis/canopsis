<template>
  <c-simple-tooltip v-if="value" :content="message" :top="top">
    <template #activator="{ props: activatorProps }">
      <v-icon
        v-bind="activatorProps"
        :color="color"
        :size="size"
        class="d-block"
      >
        sync_problem
      </v-icon>
    </template>
  </c-simple-tooltip>
</template>

<script>
import { getDiffBetweenDates } from '@/helpers/date/date';
import { convertDurationToString } from '@/helpers/date/duration';

export default {
  props: {
    value: {
      type: Number,
      default: 0,
    },
    color: {
      type: String,
      default: '',
    },
    maxWidth: {
      type: [String, Number],
      required: false,
    },
    top: {
      type: Boolean,
      required: false,
    },
    right: {
      type: Boolean,
      required: false,
    },
    bottom: {
      type: Boolean,
      required: false,
    },
    left: {
      type: Boolean,
      required: false,
    },
    size: {
      type: [String, Number],
      required: false,
    },
  },

  computed: {
    durationDiff() {
      return getDiffBetweenDates(Date.now(), this.value);
    },

    message() {
      return this.$t('common.noReceivedEvents', {
        duration: convertDurationToString(this.durationDiff),
      });
    },
  },
};
</script>
