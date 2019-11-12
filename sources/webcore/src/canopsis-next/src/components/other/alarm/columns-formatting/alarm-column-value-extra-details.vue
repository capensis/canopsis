<template lang="pug">
  v-layout
    div(v-if="alarm.v.ack")
      v-tooltip(top)
        v-icon.badge.purple.white--text(
          small,
          slot="activator"
        ) {{ $constants.EVENT_ENTITY_STYLE[$constants.EVENT_ENTITY_TYPES.ack].icon }}
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.ack') }}
          div {{ $t('common.by') }} : {{ alarm.v.ack.a }}
          div {{ $t('common.date') }} : {{ alarm.v.ack.t | date('long') }}
          div(v-if="alarm.v.ack.m") {{ $tc('common.comment') }} : {{ alarm.v.ack.m }}
    div(v-if="alarm.v.ticket")
      v-tooltip(top)
        v-icon.badge.blue.white--text(
          small,
          slot="activator"
        ) {{ $constants.EVENT_ENTITY_STYLE[$constants.EVENT_ENTITY_TYPES.declareTicket].icon }}
        div.text-md-center
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
          slot="activator"
        ) {{ $constants.EVENT_ENTITY_STYLE[$constants.EVENT_ENTITY_TYPES.delete].icon }}
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.canceled') }}
          div {{ $t('common.by') }} : {{ alarm.v.canceled.a }}
          div {{ $t('common.date') }} : {{ alarm.v.canceled.t | date('long') }}
          div(v-if="alarm.v.canceled.m") {{ $tc('common.comment') }} : {{ alarm.v.canceled.m }}
    div(v-if="alarm.v.snooze")
      v-tooltip(top)
        v-icon.badge.pink.white--text(
          small,
          slot="activator"
        ) {{ $constants.EVENT_ENTITY_STYLE[$constants.EVENT_ENTITY_TYPES.snooze].icon }}
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.snooze') }}
          div {{ $t('common.by') }} : {{ alarm.v.snooze.a }}
          div {{ $t('common.date') }} : {{ alarm.v.snooze.t | date('long') }}
          div {{ $t('common.end') }} : {{ alarm.v.snooze.val | date('long') }}
    div(v-if="pbehaviors.length")
      v-tooltip(top)
        v-icon.badge.cyan.accent-2.white--text(
          small,
          slot="activator"
        ) {{ $constants.EVENT_ENTITY_STYLE[$constants.EVENT_ENTITY_TYPES.pbehaviorAdd].icon }}
        div
          strong {{ $t('alarmList.actions.iconsTitles.pbehaviors') }}
          div(v-for="pbehavior in pbehaviors")
            div.mt-2.font-weight-bold {{ pbehavior.name }}
            div {{ $t('common.author') }}: {{ pbehavior.author }}
            div {{ $t('common.type') }}: {{ pbehavior.type_ }}
            div {{ pbehavior.tstart | date('long') }} - {{ pbehavior.tstop | date('long') }}
            div(v-if="pbehavior.rrule") {{ pbehavior.rrule }}
            div(
              v-show="pbehavior.comments && pbehavior.comments.length",
              v-for="comment in pbehavior.comments",
              :key="comment._id"
            ) {{ $tc('common.comment', pbehavior.comments.length) }}:
              div.ml-2 - {{ comment.author }}: {{ comment.message }}
            v-divider
</template>

<script>
/**
 * Component for the 'extra-details' column of the alarms list
 *
 * @module alarm
 *
 * @prop {Object} alarm - Object representing the alarm
 */
export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
  computed: {
    pbehaviors() {
      return this.alarm.pbehaviors.filter((value) => {
        const now = Date.now() / 1000;

        return value.tstart <= now && now < value.tstop;
      });
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
</style>
