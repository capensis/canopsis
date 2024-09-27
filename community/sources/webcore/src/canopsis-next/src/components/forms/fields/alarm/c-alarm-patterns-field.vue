<template>
  <pattern-editor-field
    v-field="patterns"
    :disabled="disabled"
    :readonly="readonly"
    :name="name"
    :type="$constants.PATTERN_TYPES.alarm"
    :required="required"
    :attributes="availableAlarmAttributes"
    :with-type="withType"
    :counter="counter"
  >
    <template #append-count="">
      <v-btn
        v-if="counter && counter.count"
        text
        small
        @click="showPatternAlarms"
      >
        {{ $t('common.seeAlarms') }}
      </v-btn>
    </template>
  </pattern-editor-field>
</template>

<script>
import {
  isArray,
  keyBy,
  merge,
  omit,
  map,
} from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import {
  ALARM_EVENT_INITIATORS,
  ALARM_PATTERN_FIELDS,
  BASIC_ENTITY_TYPES,
  ALARM_STATES,
  ALARM_STATUSES,
  ALARM_FIELDS_TO_LABELS_KEYS,
  MAX_LIMIT,
  PATTERN_NUMBER_OPERATORS,
  PATTERN_OPERATORS,
  PATTERN_RULE_TYPES,
  PATTERN_STRING_OPERATORS,
  PATTERN_EXISTS_OPERATORS,
} from '@/constants';

import { formGroupsToPatternRulesQuery } from '@/helpers/entities/pattern/form';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { patternCountAlarmsModalMixin } from '@/mixins/pattern/pattern-count-alarms-modal';

import PatternEditorField from '@/components/forms/fields/pattern/pattern-editor-field.vue';

const { mapActions: dynamicInfoMapActions } = createNamespacedHelpers('dynamicInfo');

export default {
  components: { PatternEditorField },
  mixins: [entitiesInfoMixin, patternCountAlarmsModalMixin],
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
    excludedAttributes: {
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
    readonly: {
      type: Boolean,
      default: false,
    },
    counter: {
      type: Object,
      required: false,
    },
  },
  data() {
    return {
      infos: [],
    };
  },
  computed: {
    stringWithOneOfOptions() {
      return {
        operators: [
          ...PATTERN_STRING_OPERATORS,

          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
      };
    },

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

    infosOptions() {
      return {
        infos: this.infos,
        type: PATTERN_RULE_TYPES.infos,
      };
    },

    ticketDataOptions() {
      return {
        type: PATTERN_RULE_TYPES.object,
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
        defaultValue: ALARM_STATUSES.ongoing,
        valueField: {
          is: 'c-alarm-status-field',
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
          is: 'c-alarm-state-field',
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

    existsOptions() {
      return {
        operators: [
          ...PATTERN_EXISTS_OPERATORS,
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
        operators: [
          PATTERN_OPERATORS.equal,
          PATTERN_OPERATORS.notEqual,
          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
        valueField: {
          is: 'c-user-picker-field',
          props: (rule) => {
            const isMultiple = isArray(rule?.value);

            return {
              multiple: isMultiple,
              deletableChips: isMultiple,
              smallChips: isMultiple,
              itemValue: 'display_name',
            };
          },
        },
      };
    },

    stringWithExistAndOneOfOptions() {
      return {
        operators: [
          ...this.stringWithOneOfOptions.operators,

          PATTERN_OPERATORS.exist,
        ],
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

    metaAlarmOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.isMetaAlarm,
          PATTERN_OPERATORS.isNotMetaAlarm,
          PATTERN_OPERATORS.ruleIs,
        ],
        valueField: {
          is: 'c-meta-alarm-rule-field',
          props: {
            required: true,
          },
        },
      };
    },

    initiatorOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.equal,
          PATTERN_OPERATORS.notEqual,
          PATTERN_OPERATORS.isOneOf,
          PATTERN_OPERATORS.isNotOneOf,
        ],
        valueField: {
          is: 'c-select-field',
          props: (rule) => {
            const isMultiple = isArray(rule?.value);

            return {
              items: Object.values(ALARM_EVENT_INITIATORS),
              multiple: isMultiple,
              deletableChips: isMultiple,
              smallChips: isMultiple,
            };
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
        {
          value: ALARM_PATTERN_FIELDS.displayName,
          options: this.stringWithOneOfOptions,
        },
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
          value: ALARM_PATTERN_FIELDS.connector,
          options: this.connectorOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.connectorName,
          options: this.connectorNameOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.resource,
          options: this.resourceOptions,
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
          options: this.stringWithExistAndOneOfOptions,
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
          options: this.stringWithExistAndOneOfOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.ackInitiator,
          options: this.initiatorOptions,
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
          value: ALARM_PATTERN_FIELDS.ticketValue,
          options: this.stringWithExistAndOneOfOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.ticketMessage,
          options: this.stringWithExistAndOneOfOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.ticketInitiator,
          options: this.initiatorOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.ticketData,
          options: this.ticketDataOptions,
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
          value: ALARM_PATTERN_FIELDS.canceledInitiator,
          options: this.initiatorOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.lastComment,
          options: this.stringWithExistAndOneOfOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.lastCommentInitiator,
          options: this.initiatorOptions,
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
        {
          value: ALARM_PATTERN_FIELDS.longOutput,
          options: this.stringWithOneOfOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.initialOutput,
          options: this.stringWithOneOfOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.initialLongOutput,
          options: this.stringWithOneOfOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.changeState,
          options: this.existsOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.totalStateChanges,
          options: this.totalStateChangesOptions,
        },
        {
          value: ALARM_PATTERN_FIELDS.meta,
          options: this.metaAlarmOptions,
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
      const mergedAttributes = omit(merge(
        {},
        this.availableAttributesByValue,
        this.externalAttributesByValue,
      ), map(this.excludedAttributes, 'value'));

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

    showPatternAlarms() {
      this.showAlarmsModalByPatterns({
        alarm_pattern: formGroupsToPatternRulesQuery(this.patterns.groups),
      });
    },

    async fetchInfos() {
      const { data: infos } = await this.fetchDynamicInfosKeysWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.infos = infos;
    },
  },
};
</script>
