<template lang="pug">
  v-layout
    extra-details-ack(v-if="alarm.v.ack", :ack="alarm.v.ack")
    extra-details-last-comment(
      v-if="alarm.v.lastComment && alarm.v.lastComment.m",
      :last-comment="alarm.v.lastComment"
    )
    div(v-if="alarm.v.ticket")
      v-tooltip(top)
        v-icon.badge.blue.white--text(
          small,
          data-test="extraDetailsOpenButton-ticket",
          slot="activator"
        ) {{ $constants.EVENT_ENTITY_STYLE[$constants.EVENT_ENTITY_TYPES.declareTicket].icon }}
        div.text-md-center(:data-test="`extraDetailsContent-${alarm._id}`")
          strong {{ $t('alarmList.actions.iconsTitles.declareTicket') }}
          div {{ $t('common.by') }} : {{ alarm.v.ticket.a }}
          div {{ $t('common.date') }} : {{ alarm.v.ticket.t | dateWithToday }}
          div(
            v-if="alarm.v.ticket.val"
          ) {{ $t('alarmList.actions.iconsFields.ticketNumber') }} : {{ alarm.v.ticket.val }}
    div(v-if="alarm.v.canceled")
      v-tooltip(top)
        v-icon.badge.blue-grey.white--text(
          small,
          data-test="extraDetailsOpenButton-canceled",
          slot="activator"
        ) {{ $constants.EVENT_ENTITY_STYLE[$constants.EVENT_ENTITY_TYPES.delete].icon }}
        div.text-md-center(:data-test="`extraDetailsContent-${alarm._id}`")
          strong {{ $t('alarmList.actions.iconsTitles.canceled') }}
          div {{ $t('common.by') }} : {{ alarm.v.canceled.a }}
          div {{ $t('common.date') }} : {{ alarm.v.canceled.t | dateWithToday }}
          div.message(v-if="alarm.v.canceled.m") {{ $tc('common.comment') }} : {{ alarm.v.canceled.m }}
    extra-details-snooze(v-if="alarm.v.snooze", :snooze="alarm.v.snooze")
    div(v-if="alarm.pbehavior")
      v-tooltip(top)
        v-icon.badge.cyan.accent-2.white--text(
          small,
          data-test="extraDetailsOpenButton-pbehaviors",
          slot="activator"
        ) {{ alarm.v.pbehavior_info.icon_name }}
        div
          strong {{ $t('alarmList.actions.iconsTitles.pbehaviors') }}
          div
            div.mt-2.font-weight-bold {{ alarm.pbehavior.name }}
            div {{ $t('common.author') }}: {{ alarm.pbehavior.author }}
            div(v-if="alarm.pbehavior.type") {{ $t('common.type') }}: {{ alarm.v.pbehavior_info.type_name }}
            div(v-if="alarm.pbehavior.reason") {{ $t('common.reason') }}: {{ alarm.pbehavior.reason.name }}
            div {{ alarm.pbehavior.tstart | dateWithToday }}
              template(v-if="alarm.pbehavior.tstop")
                |  - {{ alarm.pbehavior.tstop | dateWithToday }}
            div(v-if="alarm.pbehavior.rrule") {{ alarm.pbehavior.rrule }}
            div(
              v-for="comment in alarm.pbehavior.comments",
              :key="comment._id"
            ) {{ $tc('common.comment', alarm.pbehavior.comments.length) }}:
              div.ml-2 - {{ comment.author }}: {{ comment.message }}
            v-divider
    div(v-if="alarm.causes")
      v-tooltip(top)
        v-icon.badge.brown.darken-1.white--text(
          slot="activator",
          data-test="extraDetailsOpenButton-groupCauses",
          small
        ) {{ $constants.EVENT_ENTITY_STYLE.groupCauses.icon }}
        div.text-md-center(:data-test="`extraDetailsContent-${alarm._id}`")
          strong {{ $t('alarmList.actions.iconsTitles.grouping') }}
          v-layout(row)
            v-flex
              div {{ $tc('alarmList.actions.iconsFields.rule', causesRules.length) }}&nbsp;:
            v-flex
              div(
                v-for="(rule, index) in causesRules",
                :key="rule.id",
                :style="index | ruleStyle"
              ) &nbsp;{{ rule.name }}
          div {{ $t('alarmList.actions.iconsFields.causes') }} : {{ alarm.causes.total }}
    div(v-if="alarm.consequences")
      v-tooltip(top)
        v-icon.badge.brown.darken-1.white--text(
          slot="activator",
          data-test="extraDetailsOpenButton-groupConsequences",
          small
        ) {{ $constants.EVENT_ENTITY_STYLE.groupConsequences.icon }}
        div.text-md-center(:data-test="`extraDetailsContent-${alarm._id}`")
          strong {{ $t('alarmList.actions.iconsTitles.grouping') }}
          div {{ $t('common.title') }} : {{ alarm.rule | get('name', '') }}
          div {{ $t('alarmList.actions.iconsFields.consequences') }} : {{ alarm.consequences.total }}
</template>

<script>
import { get } from 'lodash';

import ExtraDetailsLastComment from './extra-details/extra-details-last-comment.vue';
import ExtraDetailsAck from './extra-details/extra-details-ack.vue';
import ExtraDetailsSnooze from './extra-details/extra-details-snooze.vue';

/**
 * Component for the 'extra-details' column of the alarms list
 *
 * @module alarm
 *
 * @prop {Object} alarm - Object representing the alarm
 */
export default {
  components: {
    ExtraDetailsSnooze,
    ExtraDetailsAck,
    ExtraDetailsLastComment,
  },
  filters: {
    ruleStyle(index) {
      if (index % 2 === 1) {
        return { color: '#b5b5b5' };
      }

      return {};
    },
  },
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
  computed: {
    causesRules() {
      return get(this.alarm.causes, 'rules', []);
    },
  },
};
</script>

<style lang="scss" scoped>
  .badge {
    padding: 3px 7px;
    text-align: center;
    white-space: nowrap;
    vertical-align: baseline;
    border-radius: 10px;
  }
  .message {
    max-width: 600px;
    white-space: pre-line;
  }
</style>
