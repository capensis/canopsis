<template lang="pug">
  v-layout(column)
    c-collapse-panel.mb-2(v-if="withAlarm", :title="$t('common.alarmPatterns')")
      c-alarm-patterns-field(
        v-field="value.alarm_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :readonly="readonly",
        :name="alarmFieldName",
        :check-count-name="$constants.PATTERNS_FIELDS.alarm",
        :attributes="alarmAttributes",
        with-type,
        @input="errors.remove(alarmFieldName)"
      )

    c-collapse-panel.mb-2(v-if="withEntity", :title="$t('common.entityPatterns')")
      c-entity-patterns-field(
        v-field="value.entity_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :readonly="readonly",
        :name="entityFieldName",
        :check-count-name="$constants.PATTERNS_FIELDS.entity",
        :attributes="entityAttributes",
        :entity-types="entityTypes",
        with-type,
        @input="errors.remove(entityFieldName)"
      )

    c-collapse-panel.mb-2(v-if="withPbehavior", :title="$t('common.pbehaviorPatterns')")
      c-pbehavior-patterns-field(
        v-field="value.pbehavior_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :readonly="readonly",
        :name="pbehaviorFieldName",
        :check-count-name="$constants.PATTERNS_FIELDS.pbehavior",
        with-type,
        @input="errors.remove(pbehaviorFieldName)"
      )

    c-collapse-panel.mb-2(v-if="withEvent", :title="$t('common.eventPatterns')")
      c-event-filter-patterns-field(
        v-field="value.event_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :readonly="readonly",
        :name="eventFieldName",
        @input="errors.remove(eventFieldName)"
      )

    c-collapse-panel.mb-2(v-if="withTotalEntity", :title="$t('common.totalEntityPatterns')")
      c-entity-patterns-field(
        v-field="value.total_entity_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :readonly="readonly",
        :name="totalEntityFieldName",
        with-type,
        @input="errors.remove(totalEntityFieldName)"
      )

    c-collapse-panel.mb-2(v-if="withServiceWeather", :title="$t('common.serviceWeatherPatterns')")
      c-service-weather-patterns-field(
        v-field="value.weather_service_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :name="serviceWeatherFieldName",
        @input="errors.remove(serviceWeatherFieldName)"
      )
</template>

<script>
import { PATTERNS_FIELDS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    alarmAttributes: {
      type: Array,
      required: false,
    },
    entityAttributes: {
      type: Array,
      required: false,
    },
    withAlarm: {
      type: Boolean,
      default: false,
    },
    withEvent: {
      type: Boolean,
      default: false,
    },
    withEntity: {
      type: Boolean,
      default: false,
    },
    withPbehavior: {
      type: Boolean,
      default: false,
    },
    withTotalEntity: {
      type: Boolean,
      default: false,
    },
    withServiceWeather: {
      type: Boolean,
      default: false,
    },
    entityTypes: {
      type: Array,
      required: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    someRequired: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: '',
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      activePatternTab: 0,
    };
  },
  computed: {
    hasPatterns() {
      return Object.values(PATTERNS_FIELDS).some(key => this.value[key]?.groups?.length);
    },

    isPatternRequired() {
      return this.someRequired ? !this.hasPatterns : this.required;
    },

    alarmFieldName() {
      return this.preparePatternsFieldName(PATTERNS_FIELDS.alarm);
    },

    eventFieldName() {
      return this.preparePatternsFieldName(PATTERNS_FIELDS.event);
    },

    entityFieldName() {
      return this.preparePatternsFieldName(PATTERNS_FIELDS.entity);
    },

    pbehaviorFieldName() {
      return this.preparePatternsFieldName(PATTERNS_FIELDS.pbehavior);
    },

    totalEntityFieldName() {
      return this.preparePatternsFieldName(PATTERNS_FIELDS.totalEntity);
    },

    serviceWeatherFieldName() {
      return this.preparePatternsFieldName(PATTERNS_FIELDS.serviceWeather);
    },
  },
  methods: {
    preparePatternsFieldName(name) {
      return [this.name, name].filter(Boolean).join('.');
    },
  },
};
</script>
