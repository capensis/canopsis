<template lang="pug">
  div
    v-tooltip.c-extra-details(top)
      v-icon.c-extra-details__badge.cyan.accent-2.white--text(
        small,
        slot="activator"
      ) {{ pbehavior.icon_name }}
      div
        strong {{ $t('alarmList.actions.iconsTitles.pbehaviors') }}
        div
          div.mt-2.font-weight-bold {{ pbehavior.name }}
          div {{ $t('common.author') }}: {{ pbehavior.author }}
          div {{ $t('common.type') }}: {{ pbehavior.type_name }}
          div(v-if="pbehavior.reason") {{ $t('common.reason') }}: {{ pbehavior.reason }}
          div {{ tstart }}
            template(v-if="pbehavior.tstop") &nbsp;- {{ tstop }}
          div(v-if="pbehavior.rrule") {{ pbehavior.rrule }}
          div(v-for="comment in comments", :key="comment._id")
            | {{ $tc('common.comment', comments.length) }}:
            div.ml-2 - {{ comment.author }}: {{ comment.message }}
          v-divider
</template>

<script>
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';

export default {
  props: {
    pbehavior: {
      type: Object,
      required: true,
    },
    comments: {
      type: Array,
      required: false,
    },
  },
  computed: {
    tstart() {
      return convertDateToStringWithFormatForToday(this.pbehavior.tstart);
    },

    tstop() {
      return convertDateToStringWithFormatForToday(this.pbehavior.tstop);
    },
  },
};
</script>
