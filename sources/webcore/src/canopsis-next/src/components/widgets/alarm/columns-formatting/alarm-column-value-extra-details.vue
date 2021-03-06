<template lang="pug">
  v-layout
    div(v-if="alarm.v.ack")
      v-tooltip(top)
        v-icon.badge.purple.white--text(
          small,
          data-test="extraDetailsOpenButton-ack",
          slot="activator"
        ) {{ $constants.EVENT_ENTITY_STYLE[$constants.EVENT_ENTITY_TYPES.ack].icon }}
        div.text-md-center(:data-test="`extraDetailsContent-${alarm._id}`")
          strong {{ $t('alarmList.actions.iconsTitles.ack') }}
          div {{ $t('common.by') }} : {{ alarm.v.ack.a }}
          div {{ $t('common.date') }} : {{ alarm.v.ack.t | date('long') }}
          div.message(v-if="alarm.v.ack.m") {{ $tc('common.comment') }} : {{ alarm.v.ack.m }}
    div(v-if="alarm.v.lastComment && alarm.v.lastComment.m")
      v-tooltip(top)
        v-icon.badge.white--text.purple.lighten-2(
          small,
          slot="activator"
        ) {{ $constants.EVENT_ENTITY_STYLE[$constants.EVENT_ENTITY_TYPES.comment].icon }}
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.comment') }}
          div {{ $t('common.by') }} : {{ alarm.v.lastComment.a }}
          div {{ $t('common.date') }} : {{ alarm.v.lastComment.t | date('long') }}
          div.message(v-if="alarm.v.lastComment.m") {{ $tc('common.comment') }} : {{ alarm.v.lastComment.m }}
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
          div {{ $t('common.date') }} : {{ alarm.v.ticket.t | date('long') }}
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
          div {{ $t('common.date') }} : {{ alarm.v.canceled.t | date('long') }}
          div.message(v-if="alarm.v.canceled.m") {{ $tc('common.comment') }} : {{ alarm.v.canceled.m }}
    div(v-if="alarm.v.snooze")
      v-tooltip(top)
        v-icon.badge.pink.white--text(
          small,
          data-test="extraDetailsOpenButton-snooze",
          slot="activator"
        ) {{ $constants.EVENT_ENTITY_STYLE[$constants.EVENT_ENTITY_TYPES.snooze].icon }}
        div.text-md-center(:data-test="`extraDetailsContent-${alarm._id}`")
          strong {{ $t('alarmList.actions.iconsTitles.snooze') }}
          div {{ $t('common.by') }} : {{ alarm.v.snooze.a }}
          div {{ $t('common.date') }} : {{ alarm.v.snooze.t | date('long') }}
          div {{ $t('common.end') }} : {{ alarm.v.snooze.val | date('long') }}
          div.message(v-if="alarm.v.snooze.m") {{ $tc('common.comment') }} : {{ alarm.v.snooze.m }}
    div(v-if="alarm.pbehavior")
      v-tooltip(top)
        v-icon.badge.cyan.accent-2.white--text(
          small,
          data-test="extraDetailsOpenButton-pbehaviors",
          slot="activator"
        ) {{ alarm.pbehavior.type.icon_name }}
        div(:data-test="`extraDetailsContent-${alarm._id}`")
          strong {{ $t('alarmList.actions.iconsTitles.pbehaviors') }}
          div
            div.mt-2.font-weight-bold {{ alarm.pbehavior.name }}
            div {{ $t('common.author') }}: {{ alarm.pbehavior.author }}
            div(v-if="alarm.pbehavior.type") {{ $t('common.type') }}: {{ alarm.pbehavior.type.name }}
            div(v-if="alarm.pbehavior.reason") {{ $t('common.reason') }}: {{ alarm.pbehavior.reason.name }}
            div {{ alarm.pbehavior.tstart | date('long') }}
              template(v-if="alarm.pbehavior.tstop")
                |  - {{ alarm.pbehavior.tstop | date('long') }}
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

/**
 * Component for the 'extra-details' column of the alarms list
 *
 * @module alarm
 *
 * @prop {Object} alarm - Object representing the alarm
 */
export default {
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
