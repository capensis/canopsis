<template lang="pug">
  c-pattern-editor-field(
    v-field="patterns",
    :disabled="disabled",
    :readonly="readonly",
    :name="name",
    :type="$constants.PATTERN_TYPES.alarm",
    :required="required",
    :attributes="availableAlarmAttributes",
    :with-type="withType",
    :check-count-name="checkCountName"
  )
</template>

<script>
import { keyBy, merge } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import {
  ALARM_ACK_INITIATORS,
  ALARM_PATTERN_FIELDS,
  BASIC_ENTITY_TYPES,
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  ALARM_FIELDS_TO_LABELS_KEYS,
  MAX_LIMIT,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_RULE_TYPES,
  PATTERN_STRING_OPERATORS,
} from '@/constants';

import { entitiesInfoMixin } from '@/mixins/entities/info';

const { mapActions: dynamicInfoMapActions } = createNamespacedHelpers('dynamicInfo');

export default {
  mixins: [entitiesInfoMixin],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      required: true,
    },
    attributes: {
      type: Array,
      default: () => [],
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
    withType: {
      type: Boolean,
      default: false,
    },
    checkCountName: {
      type: String,
      required: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      infos: [],
    };
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

    entitiesValueField() {
      return {
        is: 'c-entity-field',
        props: {
          required: true,
        },
      };
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

    entitiesOptions() {
      return {
        operators: this.entitiesOperators,
        defaultValue: '',
        valueField: this.entitiesValueField,
      };
    },

    infosOptions() {
      return {
        infos: this.infos,
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

    tagsOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.with,
          PATTERN_OPERATORS.without,
        ],
        valueField: {
          is: 'c-alarm-tag-field',
        },
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

    ackByOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-user-picker-field',
        },
      };
    },

    stringWithExistOptions() {
      return {
        operators: [...PATTERN_STRING_OPERATORS, PATTERN_OPERATORS.exist],
      };
    },

    activatedOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.activated,
          PATTERN_OPERATORS.inactive,
        ],
      };
    },

    ackInitiatorOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-select-field',
          props: {
            items: Object.values(ALARM_ACK_INITIATORS),
          },
        },
      };
    },

    totalStateChangesOptions() {
      return {
        type: PATTERN_RULE_TYPES.number,
        operators: PATTERN_NUMBER_OPERATORS,
      };
    },

    alarmAttributes() {
      return [
        { value: ALARM_PATTERN_FIELDS.displayName },
        {
          value: ALARM_PATTERN_FIELDS.state,
          options: this.stateOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.status,
          options: this.statusOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.component,
          options: this.componentOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.resource,
          options: this.resourceOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.connector,
          options: this.connectorOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.connectorName,
          options: this.connectorNameOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.creationDate,
          options: this.dateOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.duration,
          options: this.durationOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.infos,
          options: this.infosOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.output,
          options: this.stringWithExistOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.lastEventDate,
          options: this.dateOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.lastUpdateDate,
          options: this.dateOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.ack,
          options: this.ackOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.ackAt,
          options: this.dateOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.ackBy,
          options: this.ackByOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.ackMessage,
          options: this.stringWithExistOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.ackInitiator,
          options: this.ackInitiatorOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.resolved,
          options: this.dateOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.ticket,
          options: this.ticketOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.snooze,
          options: this.snoozeOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.canceled,
          options: this.canceledOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.lastComment,
          options: this.stringWithExistOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.tags,
          options: this.tagsOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.activated,
          options: this.activatedOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.activationDate,
          options: this.dateOptions,
        },
        { value: ALARM_PATTERN_FIELDS.longOutput },
        { value: ALARM_PATTERN_FIELDS.initialOutput },
        { value: ALARM_PATTERN_FIELDS.initialLongOutput },
        {
          value: ALARM_PATTERN_FIELDS.totalStateChanges,
          options: this.totalStateChangesOptions,
        },
      ].map(variable => ({
        ...variable,

        text: this.$tc(ALARM_FIELDS_TO_LABELS_KEYS[variable.value], 2),
      }));
    },

    availableAttributesByValue() {
      return keyBy(this.alarmAttributes, 'value');
    },

    externalAttributesByValue() {
      return keyBy(this.attributes, 'value');
    },

    availableAlarmAttributes() {
      const mergedAttributes = merge(
        {},
        this.availableAttributesByValue,
        this.externalAttributesByValue,
      );

      return Object.values(mergedAttributes);
    },
  },
  mounted() {
    if (this.isProVersion) {
      this.fetchInfos();
    }
  },
  methods: {
    ...dynamicInfoMapActions({ fetchDynamicInfosKeysWithoutStore: 'fetchInfosKeysWithoutStore' }),

    async fetchInfos() {
      const { data: infos } = await this.fetchDynamicInfosKeysWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.infos = infos;
    },
  },
};
</script>
