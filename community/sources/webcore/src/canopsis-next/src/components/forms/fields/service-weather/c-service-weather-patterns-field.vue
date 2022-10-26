<template lang="pug">
  c-pattern-editor-field(
    v-field="patterns",
    :disabled="disabled",
    :name="name",
    :required="required",
    :attributes="availableServiceWeatherAttributes",
    :with-type="withType"
  )
</template>

<script>
import { keyBy, merge } from 'lodash';

import {
  SERVICE_WEATHER_PATTERN_FIELDS,
  PATTERN_OPERATORS,
  ENTITIES_STATES,
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
  },
  computed: {
    greyOptions() {
      return {
        operators: [
          PATTERN_OPERATORS.isGrey,
          PATTERN_OPERATORS.isNotGrey,
        ],
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

    iconOptions() {
      return {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-service-weather-icon-field',
        },
      };
    },

    entityAttributes() {
      return [
        {
          text: this.$t('serviceWeather.grey'),
          value: SERVICE_WEATHER_PATTERN_FIELDS.grey,
          options: this.greyOptions,
        },
        {
          text: this.$t('serviceWeather.primaryIcon'),
          value: SERVICE_WEATHER_PATTERN_FIELDS.primaryIcon,
          options: this.iconOptions,
        },
        {
          text: this.$t('serviceWeather.secondaryIcon'),
          value: SERVICE_WEATHER_PATTERN_FIELDS.secondaryIcon,
          options: this.iconOptions,
        },
        {
          text: this.$t('common.state'),
          value: SERVICE_WEATHER_PATTERN_FIELDS.state,
          options: this.stateOptions,
        },
      ];
    },

    availableAttributesByValue() {
      return keyBy(this.entityAttributes, 'value');
    },

    externalAttributesByValue() {
      return keyBy(this.attributes, 'value');
    },

    availableServiceWeatherAttributes() {
      const mergedAttributes = merge(
        {},
        this.availableAttributesByValue,
        this.externalAttributesByValue,
      );

      return Object.values(mergedAttributes);
    },
  },
};
</script>
