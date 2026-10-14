<template>
  <c-layout class="alarm-column-value-extra-details" wrap>
    <extra-details-ack v-if="alarm.v.ack" :ack="alarm.v.ack" />
    <extra-details-last-comment
      v-if="alarm.v.last_comment && (alarm.v.last_comment.m || alarm.v.last_comment.struct_m?.length)"
      :last-comment="alarm.v.last_comment"
    />
    <extra-details-ticket
      v-if="hasTickets"
      :tickets="alarm.v.tickets"
      :failed-ticket="alarm.v.failed_ticket"
    />
    <extra-details-canceled v-if="alarm.v.canceled" :canceled="alarm.v.canceled" />
    <extra-details-snooze
      v-if="alarm.v.snooze"
      :snooze="alarm.v.snooze"
      :has-inactive-pbehavior="hasInactivePbehavior"
    />
    <extra-details-pause
      v-if="alarm.v?.pause"
      :pause="alarm.v.pause"
      :entity="alarm.entity"
    />
    <extra-details-pbehavior
      v-if="alarm.v.pbehavior_info"
      :pbehavior-info="alarm.v.pbehavior_info"
    />
    <extra-details-parents v-if="alarm.parents" :rules="alarm.meta_alarm_rules" :total="alarm.parents" />
    <extra-details-children
      v-if="alarm.children"
      :total="alarm.children"
      :opened="alarm.opened_children"
      :closed="alarm.closed_children"
      :rule="alarm.meta_alarm_rule"
    />
    <extra-details-close-delay v-if="alarm.v.close_delay" :close-delay="alarm.v.close_delay" />
    <extra-details-close-delay-value
      v-if="alarm.v.close_delay_value"
      :close-delay-value="alarm.v.close_delay_value"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { isNotActivePbehaviorType } from '@/helpers/entities/pbehavior/form';

import ExtraDetailsAck from './extra-details/extra-details-ack.vue';
import ExtraDetailsCanceled from './extra-details/extra-details-canceled.vue';
import ExtraDetailsChildren from './extra-details/extra-details-children.vue';
import ExtraDetailsCloseDelayValue from './extra-details/extra-details-close-delay-value.vue';
import ExtraDetailsCloseDelay from './extra-details/extra-details-close-delay.vue';
import ExtraDetailsLastComment from './extra-details/extra-details-last-comment.vue';
import ExtraDetailsParents from './extra-details/extra-details-parents.vue';
import ExtraDetailsPause from './extra-details/extra-details-pause.vue';
import ExtraDetailsPbehavior from './extra-details/extra-details-pbehavior.vue';
import ExtraDetailsSnooze from './extra-details/extra-details-snooze.vue';
import ExtraDetailsTicket from './extra-details/extra-details-ticket.vue';

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
    ExtraDetailsCanceled,
    ExtraDetailsChildren,
    ExtraDetailsCloseDelay,
    ExtraDetailsCloseDelayValue,
    ExtraDetailsLastComment,
    ExtraDetailsParents,
    ExtraDetailsPause,
    ExtraDetailsPbehavior,
    ExtraDetailsSnooze,
    ExtraDetailsTicket,
  },

  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const hasTickets = computed(() => !!props.alarm.v.tickets?.length);
    const hasInactivePbehavior = computed(() => isNotActivePbehaviorType(props.alarm.v.pbehavior_info?.canonical_type),
    );

    return {
      hasInactivePbehavior,
      hasTickets,
    };
  },
};
</script>

<style lang="scss">
.alarm-column-value-extra-details {
  gap: 2px;
}
</style>
