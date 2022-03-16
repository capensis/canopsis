<template lang="pug">
  c-pattern-groups-field(
    v-field="groups",
    :disabled="disabled",
    :name="name",
    :attributes="alarmAttributes"
  )
</template>

<script>
import {
  ALARM_PATTERN_FIELDS,
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  PATTERN_OPERATORS,
  PATTERN_RULE_TYPES,
} from '@/constants';

export default {
  model: {
    prop: 'groups',
    event: 'input',
  },
  props: {
    groups: {
      type: Array,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      required: false,
    },
  },
  computed: {
    entitiesOperators() {
      return [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
        PATTERN_OPERATORS.hasOneOf,
        PATTERN_OPERATORS.hasNot,
      ];
    },

    entitiesValueField() {
      return {
        is: 'c-entity-field',
        props: {
          // TODO: Should be replaced on API data
          items: [
            { name: 'Entity 1', _id: 'entity_1' },
            { name: 'Entity 2', _id: 'entity_2' },
            { name: 'Entity 3', _id: 'entity_3' },
            { name: 'Entity 4', _id: 'entity_4' },
          ],
        },
      };
    },

    entitiesOptions() {
      return {
        operators: this.entitiesOperators,
        defaultValue: '',
        valueField: this.entitiesValueField,
      };
    },

    infosOptions() {
      return {
        // TODO: Should be replaced on API data
        infos: ['infos 1', 'infos 2'],
        type: PATTERN_RULE_TYPES.infos,
      };
    },

    dateOptions() {
      return {
        type: PATTERN_RULE_TYPES.date,
      };
    },

    durationOptions() {
      return {
        type: PATTERN_RULE_TYPES.duration,
      };
    },

    statusOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.equal,
          PATTERN_OPERATORS.notEqual,
        ],
        defaultValue: ENTITIES_STATUSES.ongoing,
        valueField: {
          is: 'c-entity-status-field',
        },
      };
    },

    stateOptions() {
      return {
        defaultValue: ENTITIES_STATES.ok,
        valueField: {
          is: 'c-entity-state-field',
        },
      };
    },

    ticketOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.ticketAssociated,
          PATTERN_OPERATORS.ticketNotAssociated,
        ],
      };
    },

    ackOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.acked,
          PATTERN_OPERATORS.notAcked,
        ],
      };
    },

    canceledOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.canceled,
          PATTERN_OPERATORS.notCanceled,
        ],
      };
    },

    snoozeOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.snoozed,
          PATTERN_OPERATORS.notSnoozed,
        ],
      };
    },

    alarmAttributes() {
      return [
        {
          text: this.$t('common.displayName'),
          value: ALARM_PATTERN_FIELDS.displayName,
        },
        {
          text: this.$t('common.state'),
          value: ALARM_PATTERN_FIELDS.state,
          options: this.stateOptions,
        },
        {
          text: this.$t('common.status'),
          value: ALARM_PATTERN_FIELDS.status,
          options: this.statusOptions,
        },
        {
          text: this.$t('alarm.component'),
          value: ALARM_PATTERN_FIELDS.component,
          options: this.entitiesOptions,
        },
        {
          text: this.$t('alarm.resource'),
          value: ALARM_PATTERN_FIELDS.resource,
          options: this.entitiesOptions,
        },
        {
          text: this.$t('alarm.connector'),
          value: ALARM_PATTERN_FIELDS.connector,
          options: this.entitiesOptions,
        },
        {
          text: this.$t('alarm.connectorName'),
          value: ALARM_PATTERN_FIELDS.connectorName,
          options: this.entitiesOptions,
        },
        {
          text: this.$t('common.created'),
          value: ALARM_PATTERN_FIELDS.creationDate,
          options: this.dateOptions,
        },
        {
          text: this.$t('common.duration'),
          value: ALARM_PATTERN_FIELDS.duration,
          options: this.durationOptions,
        },
        {
          text: this.$t('common.infos'),
          value: ALARM_PATTERN_FIELDS.infos,
          options: this.infosOptions,
        },
        {
          text: this.$t('common.output'),
          value: ALARM_PATTERN_FIELDS.output,
        },
        {
          text: this.$t('common.lastEventDate'),
          value: ALARM_PATTERN_FIELDS.lastEventDate,
          options: this.dateOptions,
        },
        {
          text: this.$t('common.updated'),
          value: ALARM_PATTERN_FIELDS.lastUpdateDate,
          options: this.dateOptions,
        },
        {
          text: this.$t('alarm.acked'),
          value: ALARM_PATTERN_FIELDS.ack,
          options: this.ackOptions,
        },
        {
          text: this.$t('alarm.ackedAt'),
          value: ALARM_PATTERN_FIELDS.ackAt,
          options: this.dateOptions,
        },
        {
          text: this.$t('alarm.resolvedAt'),
          value: ALARM_PATTERN_FIELDS.resolvedAt,
          options: this.dateOptions,
        },
        {
          text: this.$t('common.ticket'),
          value: ALARM_PATTERN_FIELDS.ticket,
          options: this.ticketOptions,
        },
        {
          text: this.$t('common.snoozed'),
          value: ALARM_PATTERN_FIELDS.snooze,
          options: this.snoozeOptions,
        },
        {
          text: this.$t('common.canceled'),
          value: ALARM_PATTERN_FIELDS.canceled,
          options: this.canceledOptions,
        },
      ];
    },
  },
};
</script>
