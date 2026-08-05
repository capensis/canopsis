<template>
  <div ref="wrapperElement" class="c-patterns-field__wrapper">
    <pattern-progress
      v-if="optimizationPending || optimizationFailedReason"
      :failed-reason="optimizationFailedReason"
      @cancel="cancelOptimization"
      @close="rejectAllSuggestions"
      @try-again="tryOptimization"
    />
    <c-progress-overlay :pending="pending" />
    <c-layout class="c-patterns-field" column>
      <c-label :label="$t('common.pattern', 2)" :required="isPatternRequired" />
      <c-alert v-if="errorMessage" type="error" density="compact">
        {{ errorMessage }}
      </c-alert>
      <c-collapse-panel
        v-if="withAlarm"
        :expanded="expanded.alarm"
        :title="alarmTitle || $t('common.alarmPatterns')"
      >
        <c-alarm-patterns-field
          :patterns="value.alarm_pattern"
          :required="isPatternRequired"
          :disabled="disabled"
          :readonly="readonly"
          :name="preparedAlarmName"
          :attributes="alarmAttributes"
          :alarm-counter="alarmPatternCounter"
          with-type
          @update:patterns="updatePatternAndClear(PATTERNS_FIELDS.alarm, preparedAlarmName, $event)"
          @show:alarms="showPatternAlarmsModal([PATTERNS_FIELDS.alarm])"
        />
      </c-collapse-panel>
      <c-collapse-panel
        v-if="withEntity"
        :expanded="expanded.entity"
        :title="entityTitle || $t('common.entityPatterns')"
      >
        <pattern-field-suggestions-wrapper
          :suggestions="optimizationSuggestions"
          :patterns="value.entity_pattern"
          :entity-attributes="entityAttributes"
          :optimized-fields-regexps="optimizedFieldsRegexps"
          :entities-count="entityPatternEntitiesCounter?.count"
          @apply:suggestion="applySuggestion"
          @reject:all="rejectAllSuggestions"
          @show:entities-comparison="showEntitiesComparisonModal"
        >
          <c-entity-patterns-field
            :patterns="value.entity_pattern"
            v-bind="entityPatternsCounters"
            :required="isPatternRequired"
            :disabled="disabled"
            :readonly="readonly"
            :name="preparedEntityName"
            :attributes="entityAttributes"
            :entity-types="entityTypes"
            with-type
            @update:value="
              updatePatternAndClear(PATTERNS_FIELDS.entity, preparedEntityName, $event)
            "
            @show:alarms="showPatternAlarmsModal([PATTERNS_FIELDS.entity])"
            @show:entities="showPatternEntitiesModal([PATTERNS_FIELDS.entity])"
          />
        </pattern-field-suggestions-wrapper>
      </c-collapse-panel>
      <c-collapse-panel
        v-if="withPbehavior"
        :expanded="expanded.pbehavior"
        :title="pbehaviorTitle || $t('common.pbehaviorPatterns')"
      >
        <c-pbehavior-patterns-field
          :patterns="value.pbehavior_pattern"
          v-bind="pbehaviorPatternsCounters"
          :required="isPatternRequired"
          :disabled="disabled"
          :readonly="readonly"
          :name="preparedPbehaviorName"
          :attributes="pbehaviorAttributes"
          with-type
          @update:value="
            updatePatternAndClear(PATTERNS_FIELDS.pbehavior, preparedPbehaviorName, $event)
          "
          @show:alarms="showPatternAlarmsModal([PATTERNS_FIELDS.pbehavior])"
          @show:entities="showPatternEntitiesModal([PATTERNS_FIELDS.pbehavior])"
        />
      </c-collapse-panel>
      <c-collapse-panel
        v-if="withEvent"
        :expanded="expanded.event"
        :title="eventTitle || $t('common.eventPatterns')"
      >
        <c-event-filter-patterns-field
          :patterns="value.event_pattern"
          :required="isPatternRequired"
          :disabled="disabled"
          :readonly="readonly"
          :name="preparedEventName"
          :attributes="eventAttributes"
          :counter="counters.event_pattern"
          @update:value="updatePatternAndClear(PATTERNS_FIELDS.event, preparedEventName, $event)"
          @show:entities="showPatternEntitiesModal([PATTERNS_FIELDS.event])"
        />
      </c-collapse-panel>
      <c-collapse-panel
        v-if="withTotalEntity"
        :expanded="expanded.totalEntity"
        :title="totalEntityTitle || $t('common.totalEntityPatterns')"
      >
        <c-entity-patterns-field
          :patterns="value.total_entity_pattern"
          :required="isPatternRequired"
          :disabled="disabled"
          :readonly="readonly"
          :name="preparedTotalEntityName"
          :attributes="totalEntityAttributes"
          with-type
          @update:value="
            updatePatternAndClear(PATTERNS_FIELDS.totalEntity, preparedTotalEntityName, $event)
          "
        />
      </c-collapse-panel>
      <c-collapse-panel
        v-if="withServiceWeather"
        :expanded="expanded.serviceWeather"
        :title="serviceWeatherTitle || $t('common.serviceWeatherPatterns')"
      >
        <c-service-weather-patterns-field
          :patterns="value.weather_service_pattern"
          :required="isPatternRequired"
          :disabled="disabled"
          :name="preparedServiceWeatherName"
          :attributes="weatherServiceAttributes"
          with-type
          @update:value="
            updatePatternAndClear(
              PATTERNS_FIELDS.serviceWeather,
              preparedServiceWeatherName,
              $event,
            )
          "
        />
      </c-collapse-panel>

      <pattern-try-optimization
        v-if="mayHaveOptimizationSuggestions && !optimizationSuccessful && !optimizationPending"
        @try:optimization="tryOptimization"
      />

      <c-alert v-if="!optimizationSuggestions.length && optimizationSuccessful" type="warning">
        <span
          v-html="$t('pattern.optimizationSuggestionsWasntFound')"
          class="font-weight-regular"
        />
      </c-alert>
      <c-alert v-if="allOverLimit" type="warning">
        <span>{{ $t('pattern.errors.countOverLimit', { count: allCount }) }}</span>
      </c-alert>
      <c-layout
class="w-100"
gap="2"
justify="end"
align="center">
        <pattern-count-message v-bind="patternsCountMessageBind" />
        <template v-if="hasAllInCounter">
          <v-btn
v-if="entityCountersType"
text
small
@click="showPatternEntitiesModal()">
            {{ $t('common.seeEntities') }}
          </v-btn>
          <v-btn
v-else
text
small
@click="showPatternAlarmsModal()">
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
      </c-layout>

      <slot :counters="counters" name="additional-counters" />
    </c-layout>
  </div>
</template>

<script>
import { computed, ref, toRef, watch } from 'vue';
import { useField, useFormContext } from 'vee-validate';
import { isString } from 'lodash';

import { PATTERNS_FIELDS } from '@/constants';

import { isValidPatternRule, formGroupsToPatternRulesQuery } from '@/helpers/entities/pattern/form';
import { formFilterToPatterns } from '@/helpers/entities/filter/form';
import { useI18n } from '@/hooks/i18n';

import { usePatternCountAlarmsModal } from './hooks/pattern-count-alarms-modal';
import { usePatternCountEntitiesModal } from './hooks/pattern-count-entities-modal';
import { usePatternCounters } from './hooks/pattern-counters';
import { usePatternOptimization } from './hooks/pattern-optimization';
import { usePatternsFieldExpand } from './hooks/patterns-field-expand';
import PatternCountMessage from './pattern-count-message.vue';
import PatternTryOptimization from './pattern-try-optimization.vue';
import PatternProgress from './pattern-progress.vue';
import PatternFieldSuggestionsWrapper from './pattern-field-suggestions-wrapper.vue';

/**
 * Generates a field pattern name by combining component name and field name
 *
 * @param {string} componentName - The component name
 * @param {string} fieldName - The field name
 * @returns {string} Combined field pattern name
 */
const getFieldPatternName = (componentName, fieldName) =>
  [componentName, fieldName].filter(Boolean).join('.');

export { PATTERNS_FIELDS_TO_EXPANDED_KEYS } from './hooks/patterns-field-expand';

export default {
  components: {
    PatternCountMessage,
    PatternTryOptimization,
    PatternProgress,
    PatternFieldSuggestionsWrapper,
  },
  emits: ['update:value'],
  props: {
    value: {
      type: Object,
      default: () => ({}),
    },
    alarmAttributes: {
      type: Array,
      default: () => [],
    },
    entityAttributes: {
      type: Array,
      default: () => [],
    },
    pbehaviorAttributes: {
      type: Array,
      default: () => [],
    },
    eventAttributes: {
      type: Array,
      default: () => [],
    },
    totalEntityAttributes: {
      type: Array,
      default: () => [],
    },
    weatherServiceAttributes: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
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
    expandedEvent: {
      type: Boolean,
      default: false,
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
    const { setFieldError } = useFormContext();
    const { showPatternAlarmsModal } = usePatternCountAlarmsModal(props);
    const { showPatternEntitiesModal } = usePatternCountEntitiesModal(props);

    const wrapperElement = ref(null);
    const expanded = ref({
      alarm: false,
      entity: false,
      pbehavior: false,
      event: props.expandedEvent,
      totalEntity: false,
      serviceWeather: false,
    });

    const preparedAlarmName = computed(
      () => props.alarmName || getFieldPatternName(props.name, PATTERNS_FIELDS.alarm),
    );

    const preparedEntityName = computed(
      () => props.entityName || getFieldPatternName(props.name, PATTERNS_FIELDS.entity),
    );

    const preparedPbehaviorName = computed(
      () => props.pbehaviorName || getFieldPatternName(props.name, PATTERNS_FIELDS.pbehavior),
    );

    const preparedEventName = computed(
      () => props.eventName || getFieldPatternName(props.name, PATTERNS_FIELDS.event),
    );

    const preparedTotalEntityName = computed(
      () => props.totalEntityName || getFieldPatternName(props.name, PATTERNS_FIELDS.totalEntity),
    );

    const preparedServiceWeatherName = computed(
      () =>
        props.serviceWeatherName || getFieldPatternName(props.name, PATTERNS_FIELDS.serviceWeather),
    );

    const hasPatterns = computed(() =>
      Object.values(PATTERNS_FIELDS).some(key => props.value[key]?.groups?.length),
    );

    const isPatternRequired = computed(() =>
      props.someRequired ? !hasPatterns.value : props.required,
    );

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
    const hasError = computed(() => isPatternRequired.value && !hasPatterns.value);

    const fieldName = computed(() => props.name || 'patterns');
    const validatePatterns = () => (hasError.value ? t('pattern.errors.required') : true);
    const { errorMessage, setValue, validate, setErrors, setTouched } = useField(
      fieldName,
      validatePatterns,
      {
        validateOnValueUpdate: false,
      },
    );

    const {
      counters,
      pending: countersPending,
      checkFilter,
      hasAllInCounter,
      allOverLimit,
      allCount,
      patternsCountMessageBind,
    } = usePatternCounters({
      counterMethod: toRef(props, 'counterMethod'),
      entityCountersType: toRef(props, 'entityCountersType'),
      bothCounters: toRef(props, 'bothCounters'),
      hasError,
      patterns,
    });

    const alarmPatternCounter = computed(() =>
      props.entityCountersType
        ? counters.value?.alarm_pattern
        : counters.value?.alarms?.alarm_pattern,
    );

    const entityPatternEntitiesCounter = computed(() =>
      props.entityCountersType
        ? counters.value?.entity_pattern
        : counters.value?.entities?.entity_pattern,
    );

    const entityPatternsCounters = computed(() => {
      if (props.entityCountersType) {
        return { entityCounter: counters.value?.entity_pattern };
      }

      return {
        alarmCounter: counters.value?.alarms?.entity_pattern,
        entityCounter: counters.value?.entities?.entity_pattern,
      };
    });

    const pbehaviorPatternsCounters = computed(() => {
      const counter = props.entityCountersType
        ? counters.value?.pbehavior_pattern
        : counters.value?.alarms?.pbehavior_pattern;

      return { [props.entityCountersType ? 'entityCounter' : 'alarmCounter']: counter };
    });

    /**
     * Validates pattern rules
     *
     * @param {Array} rules - Array of pattern rule groups
     * @returns {boolean} True if all rules are valid
     */
    const isValidPatternRules = rules =>
      !!rules.length &&
      rules.every(group =>
        group.every((rule) => {
          if (!isValidPatternRule(rule)) {
            return false;
          }

          if (isString(rule.cond.value)) {
            return rule.cond.value.length > 0;
          }

          return true;
        }),
      );

    /**
     * Shows alarms modal filtered by current patterns
     */
    const showPatternAlarms = () =>
      showPatternAlarmsModal({
        alarm_pattern: formGroupsToPatternRulesQuery(props.value.alarm_pattern?.groups),
        entity_pattern: formGroupsToPatternRulesQuery(props.value.entity_pattern?.groups),
        pbehavior_pattern: formGroupsToPatternRulesQuery(props.value.pbehavior_pattern?.groups),
      });

    /**
     * Shows entities modal filtered by current patterns
     */
    const showPatternEntities = () =>
      showPatternEntitiesModal({
        search_pattern: formGroupsToPatternRulesQuery(props.value.entity_pattern.groups),
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
    } = usePatternOptimization(
      {
        wrapperElement,
        value: toRef(props, 'value'),
      },
      emit,
    );

    const updatePattern = (field, pattern) =>
      emit('update:value', {
        ...props.value,
        [field]: pattern,
      });

    const updatePatternAndClear = (field, name, pattern) => {
      setFieldError(name, undefined);
      setErrors([]);
      setTouched(false);
      updatePattern(field, pattern);
    };

    watch(
      () => props.value,
      (value) => {
        setValue(value, false);
      },
    );

    watch(hasError, (value) => {
      if (!value) {
        setErrors([]);
        setTouched(false);
      }
    });

    watch(optimizationSuggestions, (suggestions) => {
      if (suggestions.length) {
        expanded.value.entity = true;
      }
    });

    usePatternsFieldExpand({ wrapperElement, expanded });

    return {
      PATTERNS_FIELDS,
      wrapperElement,
      expanded,
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
      hasError,
      errorMessage,
      hasAllInCounter,
      patternsCountMessageBind,
      patternsFields,
      patterns,
      allOverLimit,
      allCount,
      isValidPatternRules,
      showPatternAlarms,
      showPatternEntities,
      checkFilter,
      showPatternAlarmsModal,
      showPatternEntitiesModal,
      updatePattern,
      updatePatternAndClear,
      validate,
      alarmPatternCounter,
      entityPatternEntitiesCounter,
      entityPatternsCounters,
      pbehaviorPatternsCounters,

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
