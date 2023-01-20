<template lang="pug">
  c-pattern-editor-field(
    v-field="patterns",
    :disabled="disabled",
    :readonly="readonly",
    :name="name",
    :required="required",
    :attributes="eventFilterAttributes"
  )
</template>

<script>
import {
  BASIC_ENTITY_TYPES,
  ENTITIES_STATES,
  EVENT_ENTITY_TYPES,
  EVENT_FILTER_PATTERN_FIELDS,
  EVENT_FILTER_SOURCE_TYPES,
  PATTERN_OPERATORS,
  PATTERN_RULE_TYPES,
} from '@/constants';

export default {
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
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
    required: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    entitiesOperators() {
      return [
        PATTERN_OPERATORS.equal,
        PATTERN_OPERATORS.notEqual,
        PATTERN_OPERATORS.isOneOf,
        PATTERN_OPERATORS.isNotOneOf,
        PATTERN_OPERATORS.contains,
        PATTERN_OPERATORS.notContains,
        PATTERN_OPERATORS.regexp,
      ];
    },

    connectorOptions() {
      return {
        operators: this.entitiesOperators,
        defaultValue: '',
        valueField: {
          is: 'c-entity-field',
          props: {
            required: true,
            entityTypes: [BASIC_ENTITY_TYPES.connector],
            itemText: 'connector_type',
            itemValue: 'connector_type',
          },
        },
      };
    },

    connectorNameOptions() {
      return {
        operators: this.entitiesOperators,
        defaultValue: '',
        valueField: {
          is: 'c-entity-field',
          props: {
            required: true,
            entityTypes: [BASIC_ENTITY_TYPES.connector],
            itemText: 'name',
            itemValue: 'name',
          },
        },
      };
    },

    componentOptions() {
      return {
        operators: this.entitiesOperators,
        defaultValue: '',
        valueField: {
          is: 'c-entity-field',
          props: {
            required: true,
            entityTypes: [BASIC_ENTITY_TYPES.component],
          },
        },
      };
    },

    resourceOptions() {
      return {
        operators: this.entitiesOperators,
        defaultValue: '',
        valueField: {
          is: 'c-entity-field',
          props: {
            required: true,
            entityTypes: [BASIC_ENTITY_TYPES.resource],
            itemText: 'name',
            itemValue: 'name',
          },
        },
      };
    },

    extraInfosOptions() {
      return {
        type: PATTERN_RULE_TYPES.extraInfos,
      };
    },

    eventTypes() {
      return [
        EVENT_ENTITY_TYPES.ack,
        EVENT_ENTITY_TYPES.ackRemove,
        EVENT_ENTITY_TYPES.assocTicket,
        EVENT_ENTITY_TYPES.declareTicket,
        EVENT_ENTITY_TYPES.cancel,
        EVENT_ENTITY_TYPES.uncancel,
        EVENT_ENTITY_TYPES.changeState,
        EVENT_ENTITY_TYPES.check,
        EVENT_ENTITY_TYPES.comment,
        EVENT_ENTITY_TYPES.snooze,
      ].map(value => ({
        value,
        text: this.$t(`common.entityEventTypes.${value}`),
      }));
    },

    eventTypeOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.equal,
          PATTERN_OPERATORS.notEqual,
          PATTERN_OPERATORS.contains,
          PATTERN_OPERATORS.notContains,
          PATTERN_OPERATORS.regexp,
        ],
        valueField: {
          is: 'v-combobox',
          props: {
            items: this.eventTypes,
            returnObject: false,
            combobox: true,
          },
        },
      };
    },

    sourceTypes() {
      return [
        {
          value: EVENT_FILTER_SOURCE_TYPES.component,
          text: this.$t('common.component'),
        },
        {
          value: EVENT_FILTER_SOURCE_TYPES.resource,
          text: this.$t('common.resource'),
        },
      ];
    },

    sourceTypeOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.equal,
          PATTERN_OPERATORS.notEqual,
          PATTERN_OPERATORS.contains,
          PATTERN_OPERATORS.notContains,
          PATTERN_OPERATORS.regexp,
        ],
        valueField: {
          is: 'c-select-field',
          props: {
            items: this.sourceTypes,
          },
        },
      };
    },

    stateOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.equal,
          PATTERN_OPERATORS.notEqual,
          PATTERN_OPERATORS.higher,
          PATTERN_OPERATORS.lower,
        ],
        defaultValue: ENTITIES_STATES.ok,
        valueField: {
          is: 'c-entity-state-field',
        },
      };
    },

    eventFilterAttributes() {
      return [
        {
          text: this.$t('common.eventType'),
          value: EVENT_FILTER_PATTERN_FIELDS.eventType,
          options: this.eventTypeOptions,
        },
        {
          text: this.$t('common.state'),
          value: EVENT_FILTER_PATTERN_FIELDS.state,
          options: this.stateOptions,
        },
        {
          text: this.$t('common.sourceType'),
          value: EVENT_FILTER_PATTERN_FIELDS.sourceType,
          options: this.sourceTypeOptions,
        },
        {
          text: this.$t('common.component'),
          value: EVENT_FILTER_PATTERN_FIELDS.component,
          options: this.componentOptions,
        },
        {
          text: this.$t('common.connector'),
          value: EVENT_FILTER_PATTERN_FIELDS.connector,
          options: this.connectorOptions,
        },
        {
          text: this.$t('common.connectorName'),
          value: EVENT_FILTER_PATTERN_FIELDS.connectorName,
          options: this.connectorNameOptions,
        },
        {
          text: this.$t('common.resource'),
          value: EVENT_FILTER_PATTERN_FIELDS.resource,
          options: this.resourceOptions,
        },
        {
          text: this.$t('common.output'),
          value: EVENT_FILTER_PATTERN_FIELDS.output,
        },
        {
          text: this.$tc('common.extraInfo'),
          value: EVENT_FILTER_PATTERN_FIELDS.extraInfos,
          options: this.extraInfosOptions,
        },
        {
          text: this.$t('common.longOutput'),
          value: EVENT_FILTER_PATTERN_FIELDS.longOutput,
        },
      ];
    },
  },
};
</script>
