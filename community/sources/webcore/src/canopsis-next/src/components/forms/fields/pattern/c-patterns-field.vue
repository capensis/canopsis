<template>
  <v-layout
    class="c-patterns-field"
    column
  >
    <c-collapse-panel
      v-if="withAlarm"
      :outline-color="alarmPatternOutlineColor"
      :title="$t('common.alarmPatterns')"
    >
      <c-alarm-patterns-field
        v-field="value.alarm_pattern"
        :required="isPatternRequired"
        :disabled="disabled"
        :readonly="readonly"
        :name="alarmFieldName"
        :attributes="alarmAttributes"
        :counter="counters.alarm_pattern"
        with-type
        @input="errors.remove(alarmFieldName)"
      />
    </c-collapse-panel>
    <c-collapse-panel
      v-if="withEntity"
      :outline-color="entityPatternOutlineColor"
      :title="$t('common.entityPatterns')"
    >
      <c-entity-patterns-field
        v-field="value.entity_pattern"
        :required="isPatternRequired"
        :disabled="disabled"
        :readonly="readonly"
        :name="entityFieldName"
        :attributes="entityAttributes"
        :entity-types="entityTypes"
        :counter="counters.entity_pattern"
        with-type
        @input="errors.remove(entityFieldName)"
      />
    </c-collapse-panel>
    <c-collapse-panel
      v-if="withPbehavior"
      :outline-color="pbehaviorPatternOutlineColor"
      :title="$t('common.pbehaviorPatterns')"
    >
      <c-pbehavior-patterns-field
        v-field="value.pbehavior_pattern"
        :required="isPatternRequired"
        :disabled="disabled"
        :readonly="readonly"
        :name="pbehaviorFieldName"
        :counter="counters.pbehavior_pattern"
        with-type
        @input="errors.remove(pbehaviorFieldName)"
      />
    </c-collapse-panel>
    <c-collapse-panel
      v-if="withEvent"
      :outline-color="eventPatternOutlineColor"
      :title="$t('common.eventPatterns')"
    >
      <c-event-filter-patterns-field
        v-field="value.event_pattern"
        :required="isPatternRequired"
        :disabled="disabled"
        :readonly="readonly"
        :name="eventFieldName"
        :counter="counters.event_pattern"
        @input="errors.remove(eventFieldName)"
      />
    </c-collapse-panel>
    <c-collapse-panel
      v-if="withTotalEntity"
      :outline-color="totalEntityPatternOutlineColor"
      :title="$t('common.totalEntityPatterns')"
    >
      <c-entity-patterns-field
        v-field="value.total_entity_pattern"
        :required="isPatternRequired"
        :disabled="disabled"
        :readonly="readonly"
        :name="totalEntityFieldName"
        :counter="counters.total_entity_pattern"
        with-type
        @input="errors.remove(totalEntityFieldName)"
      />
    </c-collapse-panel>
    <c-collapse-panel
      v-if="withServiceWeather"
      :outline-color="serviceWeatherPatternOutlineColor"
      :title="$t('common.serviceWeatherPatterns')"
    >
      <c-service-weather-patterns-field
        v-field="value.weather_service_pattern"
        :required="isPatternRequired"
        :disabled="disabled"
        :name="serviceWeatherFieldName"
        :counter="counters.weather_service_pattern"
        @input="errors.remove(serviceWeatherFieldName)"
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
      <pattern-count-message
        :error="hasError"
        :message="checkFilterMessages"
      />
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
        class="mr-0 ml-4"
        :disabled="!hasPatterns"
        :loading="countersPending"
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
import { PATTERNS_FIELDS } from '@/constants';

import {
  isValidPatternRule,
  formGroupsToPatternRules,
  formGroupsToPatternRulesQuery,
} from '@/helpers/entities/pattern/form';
import { formFilterToPatterns } from '@/helpers/entities/filter/form';

import { patternCountAlarmsModalMixin } from '@/mixins/pattern/pattern-count-alarms-modal';
import { patternCountEntitiesModalMixin } from '@/mixins/pattern/pattern-count-entities-modal';

import PatternCountMessage from '@/components/forms/fields/pattern/pattern-count-message.vue';

const { mapActions: mapPatternActions } = createNamespacedHelpers('pattern');

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
    bothCounters: {
      type: Boolean,
      default: false,
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

      const allCount = this.counters?.all?.count ?? 0;
      const entitiesCount = this.counters?.entities?.count ?? 0;

      if (this.entityCountersType) {
        return this.$t('pattern.entitiesCount', { entitiesCount: allCount });
      }

      if (this.bothCounters) {
        return this.$t('pattern.alarmsEntitiesCount', {
          entitiesCount,
          alarmsCount: allCount,
        });
      }

      return this.$t('pattern.alarmsCount', { alarmsCount: allCount });
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
      this.showEntitiesModalByPatterns({
        entity_pattern: formGroupsToPatternRulesQuery(this.value.entity_pattern.groups),
      });
    },

    async checkFilter() {
      try {
        this.countersPending = true;

        const method = this.entityCountersType
          ? this.checkPatternsEntitiesCount
          : this.checkPatternsAlarmsCount;

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

      if (!this.isPatternRequired && !rules.length) {
        return undefined;
      }

      return this.isValidPatternRules(rules) ? CSS_COLORS_VARS.primary : CSS_COLORS_VARS.error;
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
