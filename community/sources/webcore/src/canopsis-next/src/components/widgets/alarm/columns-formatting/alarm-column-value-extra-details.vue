<template lang="pug">
  v-layout
    extra-details-ack(v-if="alarm.v.ack", :ack="alarm.v.ack")
    extra-details-last-comment(
      v-if="alarm.v.last_comment && alarm.v.last_comment.m",
      :last-comment="alarm.v.last_comment"
    )
    extra-details-ticket(
      v-if="alarm.v.ticket",
      :ticket="alarm.v.ticket"
    )
    extra-details-canceled(
      v-if="alarm.v.canceled",
      :canceled="alarm.v.canceled"
    )
    extra-details-snooze(
      v-if="alarm.v.snooze",
      :snooze="alarm.v.snooze"
    )
    extra-details-pbehavior(
      v-if="alarm.pbehavior",
      :pbehavior="alarm.pbehavior",
      :pbehavior-info="alarm.v.pbehavior_info"
    )
    extra-details-causes(
      v-if="alarm.causes",
      :causes="alarm.causes"
    )
    extra-details-consequences(
      v-if="alarm.consequences",
      :consequences="alarm.consequences",
      :rule="alarm.rule"
    )
</template>

<script>
import ExtraDetailsAck from './extra-details/extra-details-ack.vue';
import ExtraDetailsLastComment from './extra-details/extra-details-last-comment.vue';
import ExtraDetailsTicket from './extra-details/extra-details-ticket.vue';
import ExtraDetailsCanceled from './extra-details/extra-details-canceled.vue';
import ExtraDetailsSnooze from './extra-details/extra-details-snooze.vue';
import ExtraDetailsPbehavior from './extra-details/extra-details-pbehavior.vue';
import ExtraDetailsCauses from './extra-details/extra-details-causes.vue';
import ExtraDetailsConsequences from './extra-details/extra-details-consequences.vue';

/**
 * Component for the 'extra-details' column of the alarms list
 *
 * @module alarm
 *
 * @prop {Object} alarm - Object representing the alarm
 */
export default {
  components: {
    ExtraDetailsAck,
    ExtraDetailsLastComment,
    ExtraDetailsTicket,
    ExtraDetailsCanceled,
    ExtraDetailsSnooze,
    ExtraDetailsPbehavior,
    ExtraDetailsCauses,
    ExtraDetailsConsequences,
  },
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
};
</script>
