<template>
  <div class="c-patterns-field__wrapper">
    <pattern-optimization-progress
      v-if="optimizationPending || optimizationFailedReason"
      :failed-reason="optimizationFailedReason"
      @cancel:optimization="cancelOptimization"
      @try:optimization="tryOptimization"
    />
    <v-layout
      class="c-patterns-field"
      column
    >
      <c-collapse-panel
        v-if="withAlarm"
        :outline-color="alarmPatternOutlineColor"
        :title="alarmTitle || $t('common.alarmPatterns')"
      >
        <c-alarm-patterns-field
          v-field="value.alarm_pattern"
          :required="isPatternRequired"
          :disabled="disabled"
          :readonly="readonly"
          :name="preparedAlarmName"
          :attributes="alarmAttributes"
          :counter="counters.alarm_pattern"
          with-type
          @input="errors.remove(preparedAlarmName)"
        />
      </c-collapse-panel>
      <c-collapse-panel
        v-if="withEntity"
        :outline-color="entityPatternOutlineColor"
        :title="entityTitle || $t('common.entityPatterns')"
      >
        <c-entity-patterns-field
          v-field="value.entity_pattern"
          :required="isPatternRequired"
          :disabled="disabled"
          :readonly="readonly"
          :name="preparedEntityName"
          :attributes="entityAttributes"
          :entity-types="entityTypes"
          :counter="counters.entity_pattern"
          with-type
          @input="errors.remove(preparedEntityName)"
        />
      </c-collapse-panel>
      <c-collapse-panel
        v-if="withPbehavior"
        :outline-color="pbehaviorPatternOutlineColor"
        :title="pbehaviorTitle || $t('common.pbehaviorPatterns')"
      >
        <c-pbehavior-patterns-field
          v-field="value.pbehavior_pattern"
          :required="isPatternRequired"
          :disabled="disabled"
          :readonly="readonly"
          :name="preparedPbehaviorName"
          :counter="counters.pbehavior_pattern"
          with-type
          @input="errors.remove(preparedPbehaviorName)"
        />
      </c-collapse-panel>
      <c-collapse-panel
        v-if="withEvent"
        :outline-color="eventPatternOutlineColor"
        :title="eventTitle || $t('common.eventPatterns')"
      >
        <c-event-filter-patterns-field
          v-field="value.event_pattern"
          :required="isPatternRequired"
          :disabled="disabled"
          :readonly="readonly"
          :name="preparedEventName"
          :counter="counters.event_pattern"
          @input="errors.remove(preparedEventName)"
        />
      </c-collapse-panel>
      <c-collapse-panel
        v-if="withTotalEntity"
        :outline-color="totalEntityPatternOutlineColor"
        :title="totalEntityTitle || $t('common.totalEntityPatterns')"
      >
        <c-entity-patterns-field
          v-field="value.total_entity_pattern"
          :required="isPatternRequired"
          :disabled="disabled"
          :readonly="readonly"
          :name="preparedTotalEntityName"
          :counter="counters.total_entity_pattern"
          with-type
          @input="errors.remove(preparedTotalEntityName)"
        />
      </c-collapse-panel>
      <c-collapse-panel
        v-if="withServiceWeather"
        :outline-color="serviceWeatherPatternOutlineColor"
        :title="serviceWeatherTitle || $t('common.serviceWeatherPatterns')"
      >
        <c-service-weather-patterns-field
          v-field="value.weather_service_pattern"
          :required="isPatternRequired"
          :disabled="disabled"
          :name="preparedServiceWeatherName"
          :counter="counters.weather_service_pattern"
          with-type
          @input="errors.remove(preparedServiceWeatherName)"
        />
      </c-collapse-panel>

      <pattern-try-optimization
        v-if="mayHaveOptimizationSuggestions && !optimizationSuccessful && !optimizationPending"
        @try:optimization="tryOptimization"
      />

      <pattern-suggestions
        v-if="optimizationSuggestions.length"
        :suggestions="optimizationSuggestions"
        :patterns="value.entity_pattern"
        :entity-attributes="entityAttributes"
        :optimized-fields-regexps="optimizedFieldsRegexps"
        @apply:suggestion="applySuggestion"
        @reject:all="rejectAllSuggestions"
        @show:entities-comparison="showEntitiesComparisonModal"
      />
      <c-alert v-else-if="optimizationSuccessful" type="warning">
        <span v-html="$t('pattern.optimizationSuggestionsWasntFound')" class="font-weight-regular" />
      </c-alert>
      <c-alert
        :value="allOverLimit"
        type="warning"
        transition="fade-transition"
      >
        <span>{{ $t('pattern.errors.countOverLimit', { count: allCount }) }}</span>
      </c-alert>
      <v-layout
        justify-end
        align-center
      >
        <pattern-count-message :error="hasError">
          <span v-html="checkFilterMessages" />
        </pattern-count-message>
        <template v-if="hasAllInCounter">
          <v-btn
            v-if="entityCountersType"
            text
            small
            @click="showPatternEntities"
          >
            {{ $t('common.seeEntities') }}
          </v-btn>
          <v-btn
            v-else
            text
            small
            @click="showPatternAlarms"
          >
            {{ $t('common.seeAlarms') }}
          </v-btn>
        </template>
        <v-btn
          :disabled="!hasPatterns"
          :loading="countersPending"
          class="mr-0 ml-4"
          color="primary"
          @click="checkFilter"
        >
          {{ $t('common.checkFilter') }}
        </v-btn>
      </v-layout>
    </v-layout>
  </div>
</template>

<script>
import { computed, toRef } from 'vue';
import { isString, isEmpty } from 'lodash';

import { CSS_COLORS_VARS } from '@/config';
import { PATTERNS_FIELDS, PATTERN_DURATION_FORMAT, TIME_UNITS } from '@/constants';

import { sanitizeHtml } from '@/helpers/html';
import {
  isValidPatternRule,
  formGroupsToPatternRules,
  formGroupsToPatternRulesQuery,
} from '@/helpers/entities/pattern/form';
import { formFilterToPatterns } from '@/helpers/entities/filter/form';
import { convertDurationToString } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';

import { usePatternCountAlarmsModal } from './hooks/pattern-count-alarms-modal';
import { usePatternCountEntitiesModal } from './hooks/pattern-count-entities-modal';
import { usePatternCounters } from './hooks/pattern-counters';
import { usePatternOptimization } from './hooks/pattern-optimization';
import PatternCountMessage from './pattern-count-message.vue';
import PatternSuggestions from './pattern-suggestions.vue';
import PatternTryOptimization from './pattern-try-optimization.vue';
import PatternOptimizationProgress from './pattern-optimization-progress.vue';

/**
 * Generates a field pattern name by combining component name and field name
 *
 * @param {string} componentName - The component name
 * @param {string} fieldName - The field name
 * @returns {string} Combined field pattern name
 */
const getFieldPatternName = (componentName, fieldName) => [componentName, fieldName].filter(Boolean).join('.');

export default {
  inject: ['$validator'],
  components: {
    PatternCountMessage,
    PatternSuggestions,
    PatternTryOptimization,
    PatternOptimizationProgress,
  },
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
    entityCountersType: {
      type: Boolean,
      default: false,
    },
    counterMethod: {
      type: Function,
      required: false,
    },
    bothCounters: {
      type: Boolean,
      default: false,
    },
    alarmTitle: {
      type: String,
      default: '',
    },
    entityTitle: {
      type: String,
      default: '',
    },
    pbehaviorTitle: {
      type: String,
      default: '',
    },
    eventTitle: {
      type: String,
      default: '',
    },
    totalEntityTitle: {
      type: String,
      default: '',
    },
    serviceWeatherTitle: {
      type: String,
      default: '',
    },
    alarmName: {
      type: String,
      default: '',
    },
    entityName: {
      type: String,
      default: '',
    },
    pbehaviorName: {
      type: String,
      default: '',
    },
    eventName: {
      type: String,
      default: '',
    },
    totalEntityName: {
      type: String,
      default: '',
    },
    serviceWeatherName: {
      type: String,
      default: '',
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const validator = useValidator();

    const { showAlarmsModalByPatterns } = usePatternCountAlarmsModal();
    const { showEntitiesModalByPatterns } = usePatternCountEntitiesModal();

    const preparedAlarmName = computed(() => (
      props.alarmName || getFieldPatternName(props.name, PATTERNS_FIELDS.alarm)
    ));

    const preparedEntityName = computed(() => (
      props.entityName || getFieldPatternName(props.name, PATTERNS_FIELDS.entity)
    ));

    const preparedPbehaviorName = computed(() => (
      props.pbehaviorName || getFieldPatternName(props.name, PATTERNS_FIELDS.pbehavior)
    ));

    const preparedEventName = computed(() => (
      props.eventName || getFieldPatternName(props.name, PATTERNS_FIELDS.event)
    ));

    const preparedTotalEntityName = computed(() => (
      props.totalEntityName || getFieldPatternName(props.name, PATTERNS_FIELDS.totalEntity)
    ));

    const preparedServiceWeatherName = computed(() => (
      props.serviceWeatherName || getFieldPatternName(props.name, PATTERNS_FIELDS.serviceWeather)
    ));

    const hasPatterns = computed(() => (
      Object.values(PATTERNS_FIELDS).some(key => props.value[key]?.groups?.length)
    ));

    const isPatternRequired = computed(() => (
      props.someRequired ? !hasPatterns.value : props.required
    ));

    const patternNamesToFields = computed(() => ({
      [PATTERNS_FIELDS.alarm]: preparedAlarmName.value,
      [PATTERNS_FIELDS.entity]: preparedEntityName.value,
      [PATTERNS_FIELDS.event]: preparedEventName.value,
      [PATTERNS_FIELDS.totalEntity]: preparedTotalEntityName.value,
      [PATTERNS_FIELDS.pbehavior]: preparedPbehaviorName.value,
      [PATTERNS_FIELDS.serviceWeather]: preparedServiceWeatherName.value,
    }));

    const patternsFields = computed(() => {
      const FIELDS_TO_FLAGS = {
        [PATTERNS_FIELDS.alarm]: props.withAlarm,
        [PATTERNS_FIELDS.entity]: props.withEntity,
        [PATTERNS_FIELDS.event]: props.withEvent,
        [PATTERNS_FIELDS.pbehavior]: props.withPbehavior,
        [PATTERNS_FIELDS.totalEntity]: props.withTotalEntity,
        [PATTERNS_FIELDS.serviceWeather]: props.withServiceWeather,
      };

      return Object.entries(FIELDS_TO_FLAGS)
        .filter(([, value]) => value)
        .map(([key]) => key);
    });

    const patterns = computed(() => formFilterToPatterns(props.value, patternsFields.value));

    const { counters, pending: countersPending, checkFilter } = usePatternCounters({
      counterMethod: toRef(props, 'counterMethod'),
      entityCountersType: toRef(props, 'entityCountersType'),
      patterns,
    });

    const hasError = computed(() => (isPatternRequired.value && !hasPatterns.value));
    const hasAllInCounter = computed(() => (counters.value?.all?.count > 0));
    const allOverLimit = computed(() => (counters.value?.all?.over_limit ?? false));
    const allCount = computed(() => (counters.value?.all?.count ?? 0));

    const checkFilterMessages = computed(() => {
      if (hasError.value) {
        return t('pattern.errors.required');
      }

      if (isEmpty(counters.value)) {
        return '';
      }

      const alarmsCount = counters.value?.all?.count ?? 0;
      const allDuration = convertDurationToString(
        counters.value?.all?.ms,
        PATTERN_DURATION_FORMAT,
        TIME_UNITS.millisecond,
      );
      const durationMessage = t('pattern.searchTime', { duration: allDuration });

      let message = '';

      if (props.entityCountersType) {
        const entitiesCount = counters.value?.entity_pattern?.count ?? 0;

        message = t('pattern.entitiesCount', { entitiesCount });
      } else if (props.bothCounters) {
        const entitiesCount = counters.value?.entities?.count ?? 0;

        message = t('pattern.alarmsEntitiesCount', {
          alarmsCount,
          entitiesCount,
        });
      } else {
        message = t('pattern.alarmsCount', { alarmsCount });
      }

      return sanitizeHtml(`${message} / ${durationMessage}`);
    });

    /**
     * Validates pattern rules
     *
     * @param {Array} rules - Array of pattern rule groups
     * @returns {boolean} True if all rules are valid
     */
    const isValidPatternRules = rules => (
      !!rules.length && rules.every(
        group => group.every((rule) => {
          if (!isValidPatternRule(rule)) {
            return false;
          }

          if (isString(rule.cond.value)) {
            return rule.cond.value.length > 0;
          }

          return true;
        }),
      )
    );

    /**
     * Gets the outline color for a pattern field based on validation state
     *
     * @param {string} name - Pattern field name
     * @returns {string|undefined} CSS color variable or undefined
     */
    const getPatternOutlineColor = (name) => {
      const rules = formGroupsToPatternRules(props.value[name]?.groups ?? []);
      const fieldName = patternNamesToFields.value[name];

      if (validator.errors.has(fieldName)) {
        return CSS_COLORS_VARS.error;
      }

      if (!isPatternRequired.value && !rules.length) {
        return undefined;
      }

      return isValidPatternRules(rules) ? CSS_COLORS_VARS.primary : CSS_COLORS_VARS.error;
    };

    const alarmPatternOutlineColor = computed(() => getPatternOutlineColor(PATTERNS_FIELDS.alarm));
    const entityPatternOutlineColor = computed(() => getPatternOutlineColor(PATTERNS_FIELDS.entity));
    const eventPatternOutlineColor = computed(() => getPatternOutlineColor(PATTERNS_FIELDS.event));
    const totalEntityPatternOutlineColor = computed(() => getPatternOutlineColor(PATTERNS_FIELDS.totalEntity));
    const pbehaviorPatternOutlineColor = computed(() => getPatternOutlineColor(PATTERNS_FIELDS.pbehavior));
    const serviceWeatherPatternOutlineColor = computed(() => getPatternOutlineColor(PATTERNS_FIELDS.serviceWeather));

    /**
     * Shows alarms modal filtered by current patterns
     */
    const showPatternAlarms = () => showAlarmsModalByPatterns({
      alarm_pattern: formGroupsToPatternRulesQuery(props.value.alarm_pattern?.groups),
      entity_pattern: formGroupsToPatternRulesQuery(props.value.entity_pattern?.groups),
      pbehavior_pattern: formGroupsToPatternRulesQuery(props.value.pbehavior_pattern?.groups),
    });

    /**
     * Shows entities modal filtered by current patterns
     */
    const showPatternEntities = () => showEntitiesModalByPatterns({
      entity_pattern: formGroupsToPatternRulesQuery(props.value.entity_pattern.groups),
    });

    const {
      pending: optimizationPending,
      suggestions: optimizationSuggestions,
      failedReason: optimizationFailedReason,
      successful: optimizationSuccessful,
      optimizedFieldsRegexps,
      mayHaveOptimizationSuggestions,
      tryOptimization,
      cancelOptimization,
      applySuggestion,
      rejectAllSuggestions,
      showEntitiesComparisonModal,
    } = usePatternOptimization(toRef(props, 'value'), emit);

    return {
      counters,
      countersPending,
      preparedAlarmName,
      preparedEntityName,
      preparedPbehaviorName,
      preparedEventName,
      preparedTotalEntityName,
      preparedServiceWeatherName,
      hasPatterns,
      isPatternRequired,
      mayHaveOptimizationSuggestions,
      patternNamesToFields,
      alarmPatternOutlineColor,
      entityPatternOutlineColor,
      eventPatternOutlineColor,
      totalEntityPatternOutlineColor,
      pbehaviorPatternOutlineColor,
      serviceWeatherPatternOutlineColor,
      hasError,
      hasAllInCounter,
      checkFilterMessages,
      patternsFields,
      patterns,
      allOverLimit,
      allCount,
      isValidPatternRules,
      getPatternOutlineColor,
      showPatternAlarms,
      showPatternEntities,
      checkFilter,

      optimizationPending,
      optimizationSuggestions,
      optimizationFailedReason,
      optimizationSuccessful,
      optimizedFieldsRegexps,
      tryOptimization,
      cancelOptimization,
      applySuggestion,
      rejectAllSuggestions,
      showEntitiesComparisonModal,
    };
  },
};
</script>

<style lang="scss">
.c-patterns-field {
  gap: 16px;

  &__wrapper {
    position: relative;
  }
}
</style>
