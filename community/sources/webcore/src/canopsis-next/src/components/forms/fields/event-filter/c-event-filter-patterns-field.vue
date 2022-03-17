<template lang="pug">
  c-pattern-groups-field(
    v-field="groups",
    :disabled="disabled",
    :name="name",
    :attributes="eventFilterAttributes"
  )
</template>

<script>
import { EVENT_FILTER_FILTER_FIELDS, PATTERN_OPERATORS } from '@/constants';

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
      };
    },

    entitiesOptions() {
      return {
        operators: this.entitiesOperators,
        defaultValue: '',
        valueField: this.entitiesValueField,
      };
    },

    eventFilterAttributes() {
      return [
        {
          text: this.$t('common.component'),
          value: EVENT_FILTER_FILTER_FIELDS.component,
          options: this.entitiesOptions,
        },
        {
          text: this.$t('common.connector'),
          value: EVENT_FILTER_FILTER_FIELDS.connector,
          options: this.entitiesOptions,
        },
        {
          text: this.$t('common.connectorName'),
          value: EVENT_FILTER_FILTER_FIELDS.connectorName,
          options: this.entitiesOptions,
        },
        {
          text: this.$t('common.resource'),
          value: EVENT_FILTER_FILTER_FIELDS.resource,
          options: this.entitiesOptions,
        },
        {
          text: this.$t('common.output'),
          value: EVENT_FILTER_FILTER_FIELDS.output,
        },
        {
          text: this.$tc('common.extraInfo'),
          value: EVENT_FILTER_FILTER_FIELDS.extraInfos,
        },
      ];
    },
  },
};
</script>
