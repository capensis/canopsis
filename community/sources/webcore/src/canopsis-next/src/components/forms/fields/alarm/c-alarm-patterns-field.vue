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
  ALARM_PATTERN_FIELDS,
  BASIC_ENTITY_TYPES,
  ENTITIES_STATES,
  ENTITIES_STATUSES,
  MAX_LIMIT,
  PATTERN_OPERATORS,
  PATTERN_RULE_TYPES,
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
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
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

    ackByOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-user-picker-field',
        },
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
          text: this.$t('common.component'),
          value: ALARM_PATTERN_FIELDS.component,
          options: this.componentOptions,
        },
        {
          text: this.$t('common.resource'),
          value: ALARM_PATTERN_FIELDS.resource,
          options: this.resourceOptions,
        },
        {
          text: this.$t('common.connector'),
          value: ALARM_PATTERN_FIELDS.connector,
          options: this.connectorOptions,
        },
        {
          text: this.$t('common.connectorName'),
          value: ALARM_PATTERN_FIELDS.connectorName,
          options: this.connectorNameOptions,
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
          text: this.$t('common.acked'),
          value: ALARM_PATTERN_FIELDS.ack,
          options: this.ackOptions,
        },
        {
          text: this.$t('common.ackedAt'),
          value: ALARM_PATTERN_FIELDS.ackAt,
          options: this.dateOptions,
        },
        {
          text: this.$t('common.ackedBy'),
          value: ALARM_PATTERN_FIELDS.ackBy,
          options: this.ackByOptions,
        },
        {
          text: this.$t('common.resolvedAt'),
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
        {
          text: this.$t('common.lastComment'),
          value: ALARM_PATTERN_FIELDS.lastComment,
        },
      ];
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
