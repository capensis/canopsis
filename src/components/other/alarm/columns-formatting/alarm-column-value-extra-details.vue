<template lang="pug">
  v-layout
    div(v-if="alarm.v.ack")
      v-tooltip(top)
        v-chip(small, slot="activator", color="purple")
          v-icon(color="white") check
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.ack') }}
          div {{ $t('common.by') }} : {{ alarm.v.ack.a }}
          div {{ $t('common.date') }} : {{ $d(alarm.v.ack.t, 'long') }}
          div(v-if="alarm.v.ack.m") {{ $t('common.comment') }} : {{ alarm.v.ack.m }}
    div(v-if="alarm.v.ticket")
      v-tooltip(top)
        v-chip(small, slot="activator", color="blue")
          v-icon(color="white") local_play
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.declareTicket') }}
          div {{ $t('common.by') }} : {{ alarm.v.ticket.a }}
          div {{ $t('common.date') }} : {{ $d(alarm.v.ticket.t, 'long') }}
          div(
          v-if="alarm.v.ticket.val"
          ) {{ $t('alarmList.actions.iconsFields.ticketNumber') }} : {{ alarm.v.ticket.val }}
    div(v-if="alarm.v.canceled")
      v-tooltip(top)
        v-chip(small, slot="activator", color="blue-grey")
          v-icon(color="white") delete
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.canceled') }}
          div {{ $t('common.by') }} : {{ alarm.v.canceled.a }}
          div {{ $t('common.date') }} : {{ $d(alarm.v.canceled.t, 'long') }}
          div(v-if="alarm.v.canceled.m") {{ $t('common.comment') }} : {{ alarm.v.canceled.m }}
    div(v-if="alarm.v.snooze")
      v-tooltip(top)
        v-chip(small, slot="activator", color="pink")
          v-icon(color="white") alarm
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.snooze') }}
          div {{ $t('common.by') }} : {{ alarm.v.snooze.a }}
          div {{ $t('common.date') }} : {{ $d(alarm.v.snooze.t, 'long') }}
          div {{ $t('common.end') }} : {{ $d(alarm.v.snooze.val, 'long') }}
    div(v-if="pbehaviors.length")
      v-tooltip(top)
        v-chip(small, slot="activator", color="cyan accent-2")
          v-icon(color="white") event
        div.text-md-center
          strong {{ $t('alarmList.actions.iconsTitles.pbehaviors') }}
          div(v-for="pbehavior in pbehaviors")
            div {{ pbehavior.name }}
            div {{ $d(pbehavior.tstart, 'long') }} - {{ $d(pbehavior.tstop, 'long') }}
            div(v-if="pbehavior.rrule") {{ pbehavior.rrule }}
</template>

<script>
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
