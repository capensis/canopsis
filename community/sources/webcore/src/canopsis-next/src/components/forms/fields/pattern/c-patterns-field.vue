<template>
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
        :name="alarmName"
        :attributes="alarmAttributes"
        :alarm-counter="counters.alarm_pattern"
        with-type
        @input="errors.remove(alarmName)"
        @show:alarms="showPatternAlarmsModal([PATTERNS_FIELDS.alarm])"
      />
    </c-collapse-panel>
    <c-collapse-panel
      v-if="withEntity"
      :outline-color="entityPatternOutlineColor"
      :title="entityTitle || $t('common.entityPatterns')"
    >
      <c-entity-patterns-field
        v-field="value.entity_pattern"
        v-bind="entityPatternsCounters"
        :required="isPatternRequired"
        :disabled="disabled"
        :readonly="readonly"
        :name="entityName"
        :attributes="entityAttributes"
        :entity-types="entityTypes"
        with-type
        @input="errors.remove(entityName)"
        @show:alarms="showPatternAlarmsModal([PATTERNS_FIELDS.alarm])"
        @show:entities="showPatternEntitiesModal([PATTERNS_FIELDS.entity])"
      />
    </c-collapse-panel>
    <c-collapse-panel
      v-if="withPbehavior"
      :outline-color="pbehaviorPatternOutlineColor"
      :title="pbehaviorTitle || $t('common.pbehaviorPatterns')"
    >
      <c-pbehavior-patterns-field
        v-field="value.pbehavior_pattern"
        v-bind="pbehaviorPatternsCounters"
        :required="isPatternRequired"
        :disabled="disabled"
        :readonly="readonly"
        :name="pbehaviorName"
        with-type
        @input="errors.remove(pbehaviorName)"
        @show:alarms="showPatternAlarmsModal([PATTERNS_FIELDS.pbehavior])"
        @show:entities="showPatternEntitiesModal([PATTERNS_FIELDS.pbehavior])"
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
        :name="eventName"
        @input="errors.remove(eventName)"
        @show:entities="showPatternEntitiesModal([PATTERNS_FIELDS.event])"
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
        :name="totalEntityName"
        with-type
        @input="errors.remove(totalEntityName)"
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
        :name="serviceWeatherName"
        with-type
        @input="errors.remove(serviceWeatherName)"
      />
    </c-collapse-panel>
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
          @click="showPatternEntitiesModal()"
        >
          {{ $t('common.seeEntities') }}
        </v-btn>
        <v-btn
          v-else
          text
          small
          @click="showPatternAlarmsModal()"
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
</template>

<script>
import { ref, computed } from 'vue';
import { isString, isEmpty } from 'lodash';

import { CSS_COLORS_VARS } from '@/config';
import { PATTERNS_FIELDS, PATTERN_DURATION_FORMAT, TIME_UNITS } from '@/constants';

import { sanitizeHtml } from '@/helpers/html';
import { isValidPatternRule, formGroupsToPatternRules } from '@/helpers/entities/pattern/form';
import { formFilterToPatterns } from '@/helpers/entities/filter/form';
import { convertDurationToString } from '@/helpers/date/duration';

import { useStoreModuleHooks } from '@/hooks/store';
import { useI18n } from '@/hooks/i18n';
import { useValidator } from '@/hooks/validator/validator';
import { usePendingHandler } from '@/hooks/query/pending';

import PatternCountMessage from '@/components/forms/fields/pattern/pattern-count-message.vue';

import { usePatternCountAlarmsModal } from './hooks/pattern-count-alarms-modal';
import { usePatternCountEntitiesModal } from './hooks/pattern-count-entities-modal';

/**
 * Generates a field pattern name by combining component name and field name.
 *
 * @param {string} componentName - The name of the component.
 * @param {string} fieldName - The name of the field.
 * @returns {string} The combined field pattern name, or empty string if both are empty.
 */
const getFieldPatternName = (componentName, fieldName) => [componentName, fieldName].filter(Boolean).join('.');

export default {
  inject: ['$validator'],
  components: { PatternCountMessage },
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
      default() {
        return getFieldPatternName(this.name, PATTERNS_FIELDS.alarm);
      },
    },
    entityName: {
      type: String,
      default() {
        return getFieldPatternName(this.name, PATTERNS_FIELDS.entity);
      },
    },
    pbehaviorName: {
      type: String,
      default() {
        return getFieldPatternName(this.name, PATTERNS_FIELDS.pbehavior);
      },
    },
    eventName: {
      type: String,
      default() {
        return getFieldPatternName(this.name, PATTERNS_FIELDS.event);
      },
    },
    totalEntityName: {
      type: String,
      default() {
        return getFieldPatternName(this.name, PATTERNS_FIELDS.totalEntity);
      },
    },
    serviceWeatherName: {
      type: String,
      default() {
        return getFieldPatternName(this.name, PATTERNS_FIELDS.serviceWeather);
      },
    },
  },
  setup(props) {
    const validator = useValidator();
    const { errors } = validator;
    const { t } = useI18n();
    const { useActions: usePatternActions } = useStoreModuleHooks('pattern');
    const { checkPatternsEntitiesCount, checkPatternsAlarmsCount } = usePatternActions({
      checkPatternsEntitiesCount: 'checkPatternsEntitiesCount',
      checkPatternsAlarmsCount: 'checkPatternsAlarmsCount',
    });

    const counters = ref({});

    const { showPatternAlarmsModal } = usePatternCountAlarmsModal(props);
    const { showPatternEntitiesModal } = usePatternCountEntitiesModal(props);

    const hasPatterns = computed(() => Object.values(PATTERNS_FIELDS).some(key => props.value[key]?.groups?.length));

    const isPatternRequired = computed(() => (props.someRequired ? !hasPatterns.value : props.required));

    const patternNamesToFields = computed(() => ({
      [PATTERNS_FIELDS.alarm]: props.alarmName,
      [PATTERNS_FIELDS.entity]: props.entityName,
      [PATTERNS_FIELDS.event]: props.eventName,
      [PATTERNS_FIELDS.totalEntity]: props.totalEntityName,
      [PATTERNS_FIELDS.pbehavior]: props.pbehaviorName,
      [PATTERNS_FIELDS.serviceWeather]: props.serviceWeatherName,
    }));

    /**
     * Validates if pattern rules are valid.
     *
     * @param {Array<Array<Object>>} rules - Array of rule groups, where each group is an array of rule objects.
     * @returns {boolean} True if all rules are valid, false otherwise.
     */
    const isValidPatternRules = rules => !!rules.length && rules.every(
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

    /**
     * Gets the outline color for a pattern field based on validation state and pattern rules.
     *
     * @param {string} name - The pattern field name (e.g., PATTERNS_FIELDS.alarm).
     * @returns {string|undefined} CSS color variable for error, primary, or undefined if not required and empty.
     */
    const getPatternOutlineColor = (name) => {
      const rules = formGroupsToPatternRules(props.value[name]?.groups ?? []);
      const fieldName = patternNamesToFields.value[name];

      if (errors.has(fieldName)) {
        return CSS_COLORS_VARS.error;
      }

      if (!isPatternRequired.value && !rules.length) {
        return undefined;
      }

      return isValidPatternRules(rules) ? CSS_COLORS_VARS.primary : CSS_COLORS_VARS.error;
    };

    const entityPatternsCounters = computed(() => {
      if (props.entityCountersType) {
        return { entityCounter: counters.value?.entity_pattern };
      }

      return { alarmCounter: counters.value?.entity_pattern, entityCounter: counters.value?.entities };
    });

    const pbehaviorPatternsCounters = computed(() => ({ [props.entityCountersType ? 'entityCounter' : 'alarmCounter']: counters.value?.pbehavior_pattern }));

    const alarmPatternOutlineColor = computed(() => getPatternOutlineColor(PATTERNS_FIELDS.alarm));

    const entityPatternOutlineColor = computed(() => getPatternOutlineColor(PATTERNS_FIELDS.entity));

    const eventPatternOutlineColor = computed(() => getPatternOutlineColor(PATTERNS_FIELDS.event));

    const totalEntityPatternOutlineColor = computed(() => getPatternOutlineColor(PATTERNS_FIELDS.totalEntity));

    const pbehaviorPatternOutlineColor = computed(() => getPatternOutlineColor(PATTERNS_FIELDS.pbehavior));

    const serviceWeatherPatternOutlineColor = computed(() => getPatternOutlineColor(PATTERNS_FIELDS.serviceWeather));

    const hasError = computed(() => isPatternRequired.value && !hasPatterns.value);

    const hasAllInCounter = computed(() => counters.value?.all?.count > 0);

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

    const allOverLimit = computed(() => counters.value?.all?.over_limit ?? false);

    const allCount = computed(() => counters.value?.all?.count ?? 0);

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
     * Function for checking filter patterns and updating counters.
     *
     * Determines the appropriate counter method based on props and executes it
     * to fetch pattern counts. Updates counters state with results or empty object on error.
     */
    const checkFilterHandler = async () => {
      const method = props.counterMethod ?? {
        [true]: checkPatternsAlarmsCount,
        [props.entityCountersType]: checkPatternsEntitiesCount,
      }.true;

      try {
        counters.value = await method({ data: patterns.value });
      } catch (err) {
        console.error(err);

        counters.value = {};
      }
    };

    const { pending: countersPending, handler: checkFilter } = usePendingHandler(checkFilterHandler);

    return {
      PATTERNS_FIELDS,
      counters,
      countersPending,
      hasPatterns,
      isPatternRequired,
      alarmPatternOutlineColor,
      entityPatternOutlineColor,
      eventPatternOutlineColor,
      totalEntityPatternOutlineColor,
      pbehaviorPatternOutlineColor,
      serviceWeatherPatternOutlineColor,
      hasError,
      hasAllInCounter,
      checkFilterMessages,
      entityPatternsCounters,
      pbehaviorPatternsCounters,
      allOverLimit,
      allCount,
      showPatternAlarmsModal,
      showPatternEntitiesModal,
      checkFilter,
    };
  },
};
</script>

<style lang="scss">
.c-patterns-field {
  gap: 16px;
}
</style>
