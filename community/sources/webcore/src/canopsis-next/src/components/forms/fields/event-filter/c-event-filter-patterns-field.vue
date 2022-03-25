<template lang="pug">
  c-pattern-groups-field(
    v-field="groups",
    :disabled="disabled",
    :name="name",
    :attributes="eventFilterAttributes"
  )
</template>

<script>
import { BASIC_ENTITY_TYPES, EVENT_FILTER_PATTERN_FIELDS, PATTERN_OPERATORS, PATTERN_RULE_TYPES } from '@/constants';

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

    connectorOptions() {
      return {
        operators: this.entitiesOperators,
        defaultValue: '',
        valueField: {
          is: 'c-entity-field',
          props: {
            entityTypes: [BASIC_ENTITY_TYPES.connector],
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
            entityTypes: [BASIC_ENTITY_TYPES.resource],
          },
        },
      };
    },

    extraInfosOptions() {
      return {
        type: PATTERN_RULE_TYPES.extraInfos,
      };
    },

    eventFilterAttributes() {
      return [
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
          options: this.connectorOptions,
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
      ];
    },
  },
};
</script>
