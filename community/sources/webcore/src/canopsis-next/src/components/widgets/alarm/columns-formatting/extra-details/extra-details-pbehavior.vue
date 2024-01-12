<template>
  <div>
    <v-tooltip
      class="c-extra-details"
      top
    >
      <template #activator="{ on }">
        <span
          class="c-extra-details__badge"
          v-on="on"
          :style="{ backgroundColor: color }"
        >
          <v-icon
            :color="iconColor"
            small
          >
            {{ pbehaviorInfo.icon_name }}
          </v-icon>
        </span>
      </template>
      <div>
        <strong>{{ $t('alarm.actions.iconsTitles.pbehaviors') }}</strong>
        <div>
          <div class="mt-2 font-weight-bold">
            {{ pbehavior.name }}
          </div>
          <div v-if="pbehavior.author">
            {{ $t('common.author') }}: {{ pbehavior.author.display_name }}
          </div>
          <div v-if="pbehaviorInfo.type_name">
            {{ $t('common.type') }}: {{ pbehaviorInfo.type_name }}
          </div>
          <div v-if="pbehavior.reason">
            {{ $t('common.reason') }}: {{ pbehavior.reason.name }}
          </div>
          <div>
            {{ tstart }}
            <template v-if="pbehavior.tstop">
              &nbsp;- {{ tstop }}
            </template>
          </div>
          <div v-if="pbehavior.rrule">
            {{ pbehavior.rrule }}
          </div>
          <div v-if="pbehavior.last_comment">
            {{ $t('alarm.fields.lastComment') }}:
            <div class="ml-2">
              -&nbsp;
              <template v-if="pbehavior.last_comment.author">
                {{ pbehavior.last_comment.author.display_name }}:&nbsp;
              </template>{{ pbehavior.last_comment.message }}
            </div>
          </div>
          <v-divider />
        </div>
      </div>
    </v-tooltip>
  </div>
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
