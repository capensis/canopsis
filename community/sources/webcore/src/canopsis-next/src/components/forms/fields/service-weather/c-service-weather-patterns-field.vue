<template>
  <pattern-editor-field
    :patterns="patterns"
    :disabled="disabled"
    :name="name"
    :required="required"
    :type="$constants.PATTERN_TYPES.serviceWeather"
    :attributes="availableServiceWeatherAttributes"
    :with-type="withType"
    :counter="counter"
    class="c-service-weather-patterns-field"
    @update:patterns="$emit('update:value', $event)"
  />
</template>

<script>
import { keyBy, merge } from 'lodash';
import { computed } from 'vue';

import { ALARM_STATES, PATTERN_OPERATORS, SERVICE_WEATHER_PATTERN_FIELDS } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import PatternEditorField from '@/components/forms/fields/pattern/pattern-editor-field.vue';

export default {
  components: { PatternEditorField },

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
    counter: {
      type: Object,
      required: false,
    },
  },

  emits: ['update:patterns', 'update:value'],

  setup(props) {
    const { t } = useI18n();

    const availableServiceWeatherAttributes = computed(() => {
      const greyOptions = {
        operators: [PATTERN_OPERATORS.isGrey, PATTERN_OPERATORS.isNotGrey],
      };

      const stateOptions = {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        defaultValue: ALARM_STATES.ok,
        valueField: {
          is: 'c-alarm-state-field',
        },
      };

      const iconOptions = {
        operators: [PATTERN_OPERATORS.equal, PATTERN_OPERATORS.notEqual],
        valueField: {
          is: 'c-service-weather-icon-field',
        },
      };

      const entityAttributes = [
        {
          text: t('serviceWeather.grey'),
          value: SERVICE_WEATHER_PATTERN_FIELDS.grey,
          options: greyOptions,
        },
        {
          text: t('serviceWeather.primaryIcon'),
          value: SERVICE_WEATHER_PATTERN_FIELDS.primaryIcon,
          options: iconOptions,
        },
        {
          text: t('serviceWeather.secondaryIcon'),
          value: SERVICE_WEATHER_PATTERN_FIELDS.secondaryIcon,
          options: iconOptions,
        },
        {
          text: t('common.state'),
          value: SERVICE_WEATHER_PATTERN_FIELDS.state,
          options: stateOptions,
        },
      ];

      return Object.values(
        merge({}, keyBy(entityAttributes, 'value'), keyBy(props.attributes, 'value')),
      );
    });

    return {
      availableServiceWeatherAttributes,
    };
  },
};
</script>
