<template>
  <span>{{ timeString }}</span>
</template>

<script>
import { isNumber } from 'lodash';

import { convertDateToString } from '@/helpers/date/date';

export default {
  props: {
    history: {
      type: [Object, Number],
      required: true,
    },
    deletedCountMessageKey: {
      type: String,
      required: false,
    },
    archivedCountMessageKey: {
      type: String,
      required: false,
    },
    hideDeleted: {
      type: Boolean,
      default: false,
    },
    hideArchived: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    timeString() {
      if (isNumber(this.history)) {
        return this.$t('storageSetting.history.scriptLaunched', {
          launchedAt: convertDateToString(this.history),
        });
      }

      const { time, deleted, archived } = this.history || {};

      const result = [
        this.$t('storageSetting.history.scriptLaunched', {
          launchedAt: convertDateToString(time),
        }),
      ];

      if (!this.hideDeleted && isNumber(deleted)) {
        result.push(this.$t(this.deletedCountMessageKey, {
          count: deleted,
        }));
      }

      if (!this.hideArchived && isNumber(archived)) {
        result.push(this.$t(this.archivedCountMessageKey, {
          count: archived,
        }));
      }

      return result.join(' ');
    },
  },
};
</script>
