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
        :counter="counters.alarm_pattern"
        with-type
        @input="errors.remove(alarmName)"
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
        :name="entityName"
        :attributes="entityAttributes"
        :entity-types="entityTypes"
        :counter="counters.entity_pattern"
        with-type
        @input="errors.remove(entityName)"
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
        :name="pbehaviorName"
        :counter="counters.pbehavior_pattern"
        with-type
        @input="errors.remove(pbehaviorName)"
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
        :counter="counters.event_pattern"
        @input="errors.remove(eventName)"
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
        :counter="counters.total_entity_pattern"
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
        :counter="counters.weather_service_pattern"
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
</template>

<script>
import { isString, isEmpty } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

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

import { patternCountAlarmsModalMixin } from '@/mixins/pattern/pattern-count-alarms-modal';
import { patternCountEntitiesModalMixin } from '@/mixins/pattern/pattern-count-entities-modal';

import PatternCountMessage from '@/components/forms/fields/pattern/pattern-count-message.vue';

const { mapActions: mapPatternActions } = createNamespacedHelpers('pattern');

const getFieldPatternName = (componentName, fieldName) => [componentName, fieldName].filter(Boolean).join('.');

export default {
  inject: ['$validator'],
  components: { PatternCountMessage },
  mixins: [patternCountAlarmsModalMixin, patternCountEntitiesModalMixin],
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
  data() {
    return {
      activePatternTab: 0,
      counters: {},
      countersPending: false,
    };
  },
  computed: {
    hasPatterns() {
      return Object.values(PATTERNS_FIELDS).some(key => this.value[key]?.groups?.length);
    },

    isPatternRequired() {
      return this.someRequired ? !this.hasPatterns : this.required;
    },

    patternNamesToFields() {
      return {
        [PATTERNS_FIELDS.alarm]: this.alarmName,
        [PATTERNS_FIELDS.entity]: this.entityName,
        [PATTERNS_FIELDS.event]: this.eventName,
        [PATTERNS_FIELDS.totalEntity]: this.totalEntityName,
        [PATTERNS_FIELDS.pbehavior]: this.pbehaviorName,
        [PATTERNS_FIELDS.serviceWeather]: this.serviceWeatherName,
      };
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

    hasError() {
      return this.isPatternRequired && !this.hasPatterns;
    },

    hasAllInCounter() {
      return this.counters?.all?.count > 0;
    },

    checkFilterMessages() {
      if (this.hasError) {
        return this.$t('pattern.errors.required');
      }

      if (isEmpty(this.counters)) {
        return '';
      }

      const alarmsCount = this.counters?.all?.count ?? 0;
      const allDuration = convertDurationToString(
        this.counters?.all?.ms,
        PATTERN_DURATION_FORMAT,
        TIME_UNITS.millisecond,
      );
      const durationMessage = this.$t('pattern.searchTime', { duration: allDuration });

      let message = '';

      if (this.entityCountersType) {
        const entitiesCount = this.counters?.entity_pattern?.count ?? 0;

        message = this.$t('pattern.entitiesCount', { entitiesCount });
      } else if (this.bothCounters) {
        const entitiesCount = this.counters?.entities?.count ?? 0;

        message = this.$t('pattern.alarmsEntitiesCount', {
          alarmsCount,
          entitiesCount,
        });
      } else {
        message = this.$t('pattern.alarmsCount', { alarmsCount });
      }

      return sanitizeHtml(`${message} / ${durationMessage}`);
    },

    patternsFields() {
      const FIELDS_TO_FLAGS = {
        [PATTERNS_FIELDS.alarm]: this.withAlarm,
        [PATTERNS_FIELDS.entity]: this.withEntity,
        [PATTERNS_FIELDS.event]: this.withEvent,
        [PATTERNS_FIELDS.pbehavior]: this.withPbehavior,
        [PATTERNS_FIELDS.totalEntity]: this.withTotalEntity,
        [PATTERNS_FIELDS.serviceWeather]: this.withServiceWeather,
      };

      return Object.entries(FIELDS_TO_FLAGS)
        .filter(([, value]) => value)
        .map(([key]) => key);
    },

    patterns() {
      return formFilterToPatterns(this.value, this.patternsFields);
    },

    allOverLimit() {
      return this.counters?.all?.over_limit ?? false;
    },

    allCount() {
      return this.counters?.all?.count ?? 0;
    },
  },
  methods: {
    ...mapPatternActions({
      checkPatternsEntitiesCount: 'checkPatternsEntitiesCount',
      checkPatternsAlarmsCount: 'checkPatternsAlarmsCount',
    }),

    showPatternAlarms() {
      this.showAlarmsModalByPatterns({
        alarm_pattern: formGroupsToPatternRulesQuery(this.value.alarm_pattern?.groups),
        entity_pattern: formGroupsToPatternRulesQuery(this.value.entity_pattern?.groups),
        pbehavior_pattern: formGroupsToPatternRulesQuery(this.value.pbehavior_pattern?.groups),
      });
    },

    showPatternEntities() {
      const patterns = {
        entity_pattern: formGroupsToPatternRulesQuery(this.value.entity_pattern.groups),
      };

      if (this.withPbehavior) {
        patterns.pbehavior_pattern = formGroupsToPatternRulesQuery(this.value.pbehavior_pattern.groups);
      }

      if (this.withEvent) {
        patterns.event_pattern = formGroupsToPatternRulesQuery(this.value.event_pattern.groups);
      }

      if (this.withTotalEntity) {
        patterns.total_entity_pattern = formGroupsToPatternRulesQuery(this.value.total_entity_pattern.groups);
      }

      if (this.withServiceWeather) {
        patterns.weather_service_pattern = formGroupsToPatternRulesQuery(this.value.weather_service_pattern.groups);
      }

      this.showEntitiesModalByPatterns(patterns);
    },

    async checkFilter() {
      try {
        this.countersPending = true;

        const method = this.counterMethod ?? {
          [true]: this.checkPatternsAlarmsCount,
          [this.entityCountersType]: this.checkPatternsEntitiesCount,
        }.true;

        this.counters = await method({ data: this.patterns });
      } catch (err) {
        console.error(err);

        this.counters = {};
      } finally {
        this.countersPending = false;
      }
    },

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
      const fieldName = this.patternNamesToFields[name];

      if (this.errors.has(fieldName)) {
        return CSS_COLORS_VARS.error;
      }

      if (!this.isPatternRequired && !rules.length) {
        return undefined;
      }

      return this.isValidPatternRules(rules) ? CSS_COLORS_VARS.primary : CSS_COLORS_VARS.error;
    },
  },
};
</script>

<style lang="scss">
.c-patterns-field {
  gap: 16px;
}
</style>
