<template>
  <span>{{ timeString }}</span>
</template>

<script>
import { isNumber } from 'lodash';
import { computed } from 'vue';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

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
  setup(props) {
    const { t } = useI18n();

    const timeString = computed(() => {
      if (isNumber(props.history)) {
        return t('storageSetting.history.scriptLaunched', {
          launchedAt: convertDateToString(props.history),
        });
      }

      const { time, deleted, archived } = props.history || {};

      const result = [
        t('storageSetting.history.scriptLaunched', {
          launchedAt: convertDateToString(time),
        }),
      ];

      if (!props.hideDeleted && isNumber(deleted)) {
        result.push(t(props.deletedCountMessageKey, {
          count: deleted,
        }));
      }

      if (!props.hideArchived && isNumber(archived)) {
        result.push(t(props.archivedCountMessageKey, {
          count: archived,
        }));
      }

      return result.join(' ');
    });

    return {
      timeString,
    };
  },
};
</script>
