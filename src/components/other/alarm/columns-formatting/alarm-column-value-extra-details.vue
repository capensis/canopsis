<template lang="pug">
  v-layout
    div(v-if="alarm.v.ack")
      v-tooltip(top)
        v-icon.badge.purple.white--text(small, slot="activator") check
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.ack') }}
          div {{ $t('common.by') }} : {{ alarm.v.ack.a }}
          div {{ $t('common.date') }} : {{ alarm.v.ack.t | date('datetime') }}
          div(v-if="alarm.v.ack.m") {{ $t('common.comment') }} : {{ alarm.v.ack.m }}
    div(v-if="alarm.v.ticket")
      v-tooltip(top)
        v-icon.badge.blue.white--text(small, slot="activator") local_play
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.declareTicket') }}
          div {{ $t('common.by') }} : {{ alarm.v.ticket.a }}
          div {{ $t('common.date') }} : {{ alarm.v.ticket.t | date('datetime') }}
          div(
          v-if="alarm.v.ticket.val"
          ) {{ $t('alarmList.actions.iconsFields.ticketNumber') }} : {{ alarm.v.ticket.val }}
    div(v-if="alarm.v.canceled")
      v-tooltip(top)
        v-icon.badge.blue-grey.white--text(small, slot="activator") delete
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.canceled') }}
          div {{ $t('common.by') }} : {{ alarm.v.canceled.a }}
          div {{ $t('common.date') }} : {{ alarm.v.canceled.t | date('datetime') }}
          div(v-if="alarm.v.canceled.m") {{ $t('common.comment') }} : {{ alarm.v.canceled.m }}
    div(v-if="alarm.v.snooze")
      v-tooltip(top)
        v-icon.badge.pink.white--text(small, slot="activator") alarm
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.snooze') }}
          div {{ $t('common.by') }} : {{ alarm.v.snooze.a }}
          div {{ $t('common.date') }} : {{ alarm.v.snooze.t | date('datetime') }}
          div {{ $t('common.end') }} : {{ alarm.v.snooze.val | date('datetime') }}
    div(v-if="pbehaviors.length")
      v-tooltip(top)
        v-icon.badge.cyan.accent-2.white--text(small, slot="activator") event
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.pbehaviors') }}
          div(v-for="pbehavior in pbehaviors")
            div {{ pbehavior.name }}
            div {{ pbehavior.tstart | date('datetime') }} - {{ pbehavior.tstop | date('datetime') }}
            div(v-if="pbehavior.rrule") {{ pbehavior.rrule }}
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
        const now = Date.now();

        return value.tstart <= now && now < value.tstop;
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .badge {
    min-width: 10px;
    padding: 3px 7px;
    font-weight: 700;
    line-height: 1;
    text-align: center;
    white-space: nowrap;
    vertical-align: baseline;
    border-radius: 10px;
  }
</style>
