<template>
  <pattern-editor-field
    v-field="patterns"
    :disabled="disabled"
    :readonly="readonly"
    :name="name"
    :required="required"
    :attributes="eventFilterAttributes"
    :counter="counter"
  />
</template>

<script>
import { isArray } from 'lodash';

import {
  ALARM_STATES,
  ALARM_EVENT_INITIATORS,
  BASIC_ENTITY_TYPES,
  EVENT_FILTER_PATTERN_FIELDS,
  EVENT_FILTER_SOURCE_TYPES,
  PATTERN_OPERATORS,
  PATTERN_RULE_TYPES,
  PATTERN_STRING_OPERATORS,
} from '@/constants';

import CAlarmOldStateField from '@/components/forms/fields/alarm/c-alarm-old-state-field.vue';
import PatternEditorField from '@/components/forms/fields/pattern/pattern-editor-field.vue';

export default {
  components: { PatternEditorField },
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
    counter: {
      type: Object,
      required: false,
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

    eventTypeOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.equal,
          PATTERN_OPERATORS.notEqual,
          PATTERN_OPERATORS.contains,
          PATTERN_OPERATORS.notContains,
          PATTERN_OPERATORS.regexp,
          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
        valueField: {
          is: 'c-event-type-field',
          props: (rule) => {
            const isMultiple = isArray(rule?.value);

            return {
              multiple: isMultiple,
              deletableChips: isMultiple,
              smallChips: isMultiple,
            };
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
          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
        valueField: {
          is: 'c-select-field',
          props: (rule) => {
            const isMultiple = isArray(rule?.value);

            return {
              items: this.sourceTypes,
              multiple: isMultiple,
              deletableChips: isMultiple,
              smallChips: isMultiple,
            };
          },
        },
      };
    },

    initiatorOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.equal,
          PATTERN_OPERATORS.notEqual,
          PATTERN_OPERATORS.contains,
          PATTERN_OPERATORS.notContains,
          PATTERN_OPERATORS.regexp,
          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
        valueField: {
          is: 'c-select-field',
          props: {
            items: Object.values(ALARM_EVENT_INITIATORS),
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
        defaultValue: ALARM_STATES.ok,
        valueField: {
          is: CAlarmOldStateField,
        },
      };
    },

    stringWithOneOfOptions() {
      return {
        operators: [
          ...PATTERN_STRING_OPERATORS,

          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
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
          options: this.stringWithOneOfOptions,
        },
        {
          text: this.$tc('common.extraInfo'),
          value: EVENT_FILTER_PATTERN_FIELDS.extraInfos,
          options: this.extraInfosOptions,
        },
        {
          text: this.$t('common.longOutput'),
          value: EVENT_FILTER_PATTERN_FIELDS.longOutput,
          options: this.stringWithOneOfOptions,
        },
        {
          text: this.$t('common.author'),
          value: EVENT_FILTER_PATTERN_FIELDS.author,
          options: this.stringWithOneOfOptions,
        },
        {
          text: this.$t('common.initiator'),
          value: EVENT_FILTER_PATTERN_FIELDS.initiator,
          options: this.initiatorOptions,
        },
      ];
    },
  },
};
</script>
