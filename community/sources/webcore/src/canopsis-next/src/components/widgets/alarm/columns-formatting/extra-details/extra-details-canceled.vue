<template lang="pug">
  div
    v-tooltip.c-extra-details(top)
      v-icon.c-extra-details__badge.blue-grey.white--text(
        small,
        slot="activator"
      ) {{ icon }}
      div.text-md-center
        strong {{ $t('alarmList.actions.iconsTitles.canceled') }}
        div {{ $t('common.by') }} : {{ canceled.a }}
        div {{ $t('common.date') }} : {{ date }}
        div.c-extra-details__message(
          v-if="canceled.m"
        ) {{ $tc('common.comment') }} : {{ canceled.m }}
</template>

<script>
import { EVENT_ENTITY_TYPES } from '@/constants';

import { getEntityEventIcon } from '@/helpers/icon';
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';

export default {
  props: {
    canceled: {
      type: Object,
      required: true,
    },
  },
  computed: {
    date() {
      return convertDateToStringWithFormatForToday(this.canceled.t);
    },

    icon() {
      return getEntityEventIcon(EVENT_ENTITY_TYPES.delete);
    },
  },
};
</script>
