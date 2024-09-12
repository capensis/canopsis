<template lang="pug">
  div
    v-tooltip.c-extra-details(top, lazy)
      template(#activator="{ on }")
        span.c-extra-details__badge(v-on="on", :style="{ backgroundColor: color }")
          v-icon(:color="iconColor", small) {{ pbehaviorInfo.icon_name }}
      div
        strong {{ $t('alarm.actions.iconsTitles.pbehaviors') }}
        div
          div.mt-2.font-weight-bold {{ pbehavior.name }}
          div(v-if="pbehavior.author") {{ $t('common.author') }}: {{ pbehavior.author.display_name }}
          div(v-if="pbehaviorInfo.type_name") {{ $t('common.type') }}: {{ pbehaviorInfo.type_name }}
          div(v-if="pbehavior.reason") {{ $t('common.reason') }}: {{ pbehavior.reason.name }}
          div {{ tstart }}
            template(v-if="pbehavior.tstop") &nbsp;- {{ tstop }}
          div(v-if="pbehavior.rrule") {{ pbehavior.rrule }}
          div(v-if="pbehavior.last_comment") {{ $t('alarm.fields.lastComment') }}:
            div.ml-2 -&nbsp;
              template(v-if="pbehavior.last_comment.author") {{ pbehavior.last_comment.author.display_name }}:&nbsp;
              | {{ pbehavior.last_comment.message }}
          v-divider
</template>

<script>
import { convertDateToStringWithFormatForToday } from '@/helpers/date/date';
import { getMostReadableTextColor } from '@/helpers/color';
import { getPbehaviorColor } from '@/helpers/entities/pbehavior/form';

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
    color() {
      return getPbehaviorColor(this.pbehavior);
    },

    iconColor() {
      return getMostReadableTextColor(this.color, { level: 'AA', size: 'large' });
    },

    tstart() {
      return convertDateToStringWithFormatForToday(this.pbehavior.tstart);
    },

    tstop() {
      return convertDateToStringWithFormatForToday(this.pbehavior.tstop);
    },
  },
};
</script>
