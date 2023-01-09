<template lang="pug">
  div
    v-tooltip.c-extra-details(top)
      v-icon.c-extra-details__badge.secondary.accent-2.white--text(
        small,
        slot="activator"
      ) {{ pbehaviorInfo.icon_name }}
      div
        strong {{ $t('alarmList.actions.iconsTitles.pbehaviors') }}
        div
          div.mt-2.font-weight-bold {{ pbehavior.name }}
          div {{ $t('common.author') }}: {{ pbehavior.author.name }}
          div(v-if="pbehaviorInfo.type_name") {{ $t('common.type') }}: {{ pbehaviorInfo.type_name }}
          div(v-if="pbehavior.reason") {{ $t('common.reason') }}: {{ pbehavior.reason.name }}
          div {{ tstart }}
            template(v-if="pbehavior.tstop") &nbsp;- {{ tstop }}
          div(v-if="pbehavior.rrule") {{ pbehavior.rrule }}
          div(v-if="pbehavior.last_comment") {{ $t('common.lastComment') }}:
            div.ml-2 - {{ pbehavior.last_comment.author.name }}: {{ pbehavior.last_comment.message }}
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
    pbehaviorInfo: {
      type: Object,
      default: () => ({}),
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
