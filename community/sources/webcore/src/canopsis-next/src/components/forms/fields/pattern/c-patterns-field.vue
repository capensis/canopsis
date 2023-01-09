<template lang="pug">
  v-layout.c-patterns-field(column)
    c-collapse-panel(
      v-if="withAlarm",
      :outline-color="alarmPatternOutlineColor",
      :title="$t('common.alarmPatterns')"
    )
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

    c-collapse-panel(
      v-if="withEntity",
      :outline-color="entityPatternOutlineColor",
      :title="$t('common.entityPatterns')"
    )
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

    c-collapse-panel(
      v-if="withPbehavior",
      :outline-color="pbehaviorPatternOutlineColor",
      :title="$t('common.pbehaviorPatterns')"
    )
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

    c-collapse-panel(
      v-if="withEvent",
      :outline-color="eventPatternOutlineColor",
      :title="$t('common.eventPatterns')"
    )
      c-event-filter-patterns-field(
        v-field="value.event_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :readonly="readonly",
        :name="eventFieldName",
        @input="errors.remove(eventFieldName)"
      )

    c-collapse-panel(
      v-if="withTotalEntity",
      :outline-color="totalEntityPatternOutlineColor",
      :title="$t('common.totalEntityPatterns')"
    )
      c-entity-patterns-field(
        v-field="value.total_entity_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :readonly="readonly",
        :name="totalEntityFieldName",
        with-type,
        @input="errors.remove(totalEntityFieldName)"
      )

    c-collapse-panel(
      v-if="withServiceWeather",
      :outline-color="serviceWeatherPatternOutlineColor",
      :title="$t('common.serviceWeatherPatterns')"
    )
      c-service-weather-patterns-field(
        v-field="value.weather_service_pattern",
        :required="isPatternRequired",
        :disabled="disabled",
        :name="serviceWeatherFieldName",
        @input="errors.remove(serviceWeatherFieldName)"
      )

    v-messages(v-if="someRequired && !hasPatterns", :value="[$t('pattern.errors.required')]", color="error")
</template>

<script>
import { isString } from 'lodash';

import { PATTERNS_FIELDS } from '@/constants';

import { COLORS } from '@/config';

import { isValidPatternRule } from '@/helpers/pattern';
import { formGroupsToPatternRules } from '@/helpers/forms/pattern';

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

    alarmPatternOutlineColor() {
      return this.getPatternOutlineColor(PATTERNS_FIELDS.alarm);
    },

    entityPatternOutlineColor() {
      return this.getPatternOutlineColor(PATTERNS_FIELDS.entity);
    },

    eventPatternOutlineColor() {
      return this.getPatternOutlineColor(PATTERNS_FIELDS.event);
    },

    totalEntityPatternOutlineColor() {
      return this.getPatternOutlineColor(PATTERNS_FIELDS.totalEntity);
    },

    pbehaviorPatternOutlineColor() {
      return this.getPatternOutlineColor(PATTERNS_FIELDS.pbehavior);
    },

    serviceWeatherPatternOutlineColor() {
      return this.getPatternOutlineColor(PATTERNS_FIELDS.serviceWeather);
    },
  },
  methods: {
    isValidPatternRules(rules) {
      return !!rules.length && rules.every(
        group => group.every((rule) => {
          if (!isValidPatternRule(rule)) {
            return false;
          }

          if (isString(rule.cond.value)) {
            return rule.cond.value.length > 0;
          }

          return true;
        }),
      );
    },

    getPatternOutlineColor(name) {
      const rules = formGroupsToPatternRules(this.value[name]?.groups ?? []);

      if (!this.isPatternRequired && !rules.length) {
        return undefined;
      }

      return this.isValidPatternRules(rules) ? COLORS.primary : COLORS.error;
    },

    preparePatternsFieldName(name) {
      return [this.name, name].filter(Boolean).join('.');
    },
  },
};
</script>

<style lang="scss">
.c-patterns-field {
  gap: 16px;
}
</style>
