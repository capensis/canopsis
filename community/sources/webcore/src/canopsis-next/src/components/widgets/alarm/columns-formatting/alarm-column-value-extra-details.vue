<template>
  <v-layout
    class="alarm-column-value-extra-details"
    wrap
  >
    <extra-details-ack
      v-if="alarm.v.ack"
      :ack="alarm.v.ack"
    />
    <extra-details-last-comment
      v-if="alarm.v.last_comment && alarm.v.last_comment.m"
      :last-comment="alarm.v.last_comment"
    />
    <extra-details-ticket
      v-if="hasTickets"
      :tickets="alarm.v.tickets"
    />
    <extra-details-canceled
      v-if="alarm.v.canceled"
      :canceled="alarm.v.canceled"
    />
    <extra-details-snooze
      v-if="alarm.v.snooze"
      :snooze="alarm.v.snooze"
    />
    <extra-details-pbehavior
      v-if="alarm.pbehavior"
      :pbehavior="alarm.pbehavior"
      :pbehavior-info="alarm.v.pbehavior_info"
    />
    <extra-details-parents
      v-if="alarm.parents"
      :rules="alarm.meta_alarm_rules"
      :total="alarm.parents"
    />
    <extra-details-children
      v-if="alarm.children"
      :total="alarm.children"
      :opened="alarm.opened_children"
      :closed="alarm.closed_children"
      :rule="alarm.meta_alarm_rule"
    />
  </v-layout>
</template>

<script>
import ExtraDetailsAck from './extra-details/extra-details-ack.vue';
import ExtraDetailsLastComment from './extra-details/extra-details-last-comment.vue';
import ExtraDetailsTicket from './extra-details/extra-details-ticket.vue';
import ExtraDetailsCanceled from './extra-details/extra-details-canceled.vue';
import ExtraDetailsSnooze from './extra-details/extra-details-snooze.vue';
import ExtraDetailsPbehavior from './extra-details/extra-details-pbehavior.vue';
import ExtraDetailsParents from './extra-details/extra-details-parents.vue';
import ExtraDetailsChildren from './extra-details/extra-details-children.vue';

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
    ExtraDetailsParents,
    ExtraDetailsChildren,
  },
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
  computed: {
    hasTickets() {
      return this.alarm.v.tickets?.length;
    },
  },
};
</script>

<style lang="scss">
.alarm-column-value-extra-details {
  gap: 2px;
}
</style>
