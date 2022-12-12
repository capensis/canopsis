<template lang="pug">
  v-tooltip(top)
    template(#activator="{ on }")
      v-icon.c-extra-details__badge.blue-grey.white--text(v-on="on", small) {{ icon }}
    div.text-md-center
      strong {{ $t('alarm.actions.iconsTitles.canceled') }}
      div {{ $t('common.by') }} : {{ canceled.a }}
      div {{ $t('common.date') }} : {{ date }}
      div.c-extra-details__message(
        v-if="canceled.m"
      ) {{ $tc('common.comment') }} : {{ canceled.m }}
</template>

<script>
import { EVENT_ENTITY_STYLE, EVENT_ENTITY_TYPES } from '@/constants';

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
      return EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.delete].icon;
    },
  },
};
</script>
