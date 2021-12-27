<template lang="pug">
  v-layout(row)
    v-flex(xs12)
      v-tabs(
        v-model="activePatternTab",
        :slider-color="hasAnyError ? 'error' : 'primary'",
        color="transparent",
        fixed-tabs
      )
        template(v-if="alarm")
          v-tab(:class="{ 'error--text': errors.has(alarmPatternsFieldName) }") {{ $t('common.alarmPatterns') }}
          v-tab-item
            patterns-list(
              v-field="value.alarm_patterns",
              v-validate="rules",
              :disabled="disabled",
              :name="alarmPatternsFieldName",
              :type="$constants.PATTERNS_TYPES.alarm",
              @input="errors.remove(alarmPatternsFieldName)"
            )
        template(v-if="event")
          v-tab(:class="{ 'error--text': errors.has(eventPatternsFieldName) }") {{ $t('common.eventPatterns') }}
          v-tab-item
            patterns-list(
              v-field="value.event_patterns",
              v-validate="rules",
              :disabled="disabled",
              :name="eventPatternsFieldName",
              :type="$constants.PATTERNS_TYPES.event",
              @input="errors.remove(eventPatternsFieldName)"
            )
        template(v-if="entity")
          v-tab(:class="{ 'error--text': errors.has(entityPatternsFieldName) }") {{ $t('common.entityPatterns') }}
          v-tab-item
            patterns-list(
              v-field="value.entity_patterns",
              v-validate="rules",
              :disabled="disabled",
              :name="entityPatternsFieldName",
              :type="$constants.PATTERNS_TYPES.entity",
              @input="errors.remove(entityPatternsFieldName)"
            )
        template(v-if="totalEntity")
          v-tab(
            :class="{ 'error--text': errors.has(totalEntityPatternsFieldName) }"
          ) {{ $t('common.totalEntityPatterns') }}
          v-tab-item
            patterns-list(
              v-field="value.total_entity_patterns",
              v-validate="rules",
              :disabled="disabled",
              :name="totalEntityPatternsFieldName",
              :type="$constants.PATTERNS_TYPES.totalEntity",
              @input="errors.remove(totalEntityPatternsFieldName)"
            )
        slot(name="additionalTabs")
</template>

<script>
import { formValidationHeaderMixin } from '@/mixins/form';

import PatternsList from '@/components/common/patterns-list/patterns-list.vue';

export default {
  inject: ['$validator'],
  components: { PatternsList },
  mixins: [formValidationHeaderMixin],
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
    alarm: {
      type: Boolean,
      default: false,
    },
    event: {
      type: Boolean,
      default: false,
    },
    entity: {
      type: Boolean,
      default: false,
    },
    totalEntity: {
      type: Boolean,
      default: false,
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
  },
  data() {
    return {
      activePatternTab: 0,
    };
  },
  computed: {
    hasPatterns() {
      return Object.values(this.value).some((patterns = []) => patterns.length);
    },

    rules() {
      return {
        required: this.someRequired ? !this.hasPatterns : this.required,
      };
    },

    alarmPatternsFieldName() {
      return this.preparePatternsFieldName('alarm_patterns');
    },

    eventPatternsFieldName() {
      return this.preparePatternsFieldName('event_patterns');
    },

    entityPatternsFieldName() {
      return this.preparePatternsFieldName('entity_patterns');
    },

    totalEntityPatternsFieldName() {
      return this.preparePatternsFieldName('total_entity_patterns');
    },
  },
  methods: {
    preparePatternsFieldName(name) {
      return [this.name, name].filter(Boolean).join('.');
    },
  },
};
</script>
