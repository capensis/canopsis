<template lang="pug">
  v-layout(column)
    c-collapse-panel.mb-2(v-if="withAlarm", :error="errors.has(alarmPatternsFieldName)", color="grey")
      template(#header="")
        span.white--text {{ $t('common.alarmPatterns') }}
      v-card
        v-card-text
          c-alarm-patterns-field(
            v-field="value.alarm_pattern",
            :required="isPatternRequired",
            :disabled="disabled",
            :name="alarmPatternsFieldName",
            with-type,
            @input="errors.remove(alarmPatternsFieldName)"
          )

    c-collapse-panel.mb-2(v-if="withEntity", :error="errors.has(entityPatternsFieldName)", color="grey")
      template(#header="")
        span.white--text {{ $t('common.entityPatterns') }}
      v-card
        v-card-text
          c-entity-patterns-field(
            v-field="value.entity_pattern",
            :required="isPatternRequired",
            :disabled="disabled",
            :name="entityPatternsFieldName",
            with-type,
            @input="errors.remove(entityPatternsFieldName)"
          )

    c-collapse-panel(v-if="withPbehavior", :error="errors.has(entityPbehaviorFieldName)", color="grey")
      template(#header="")
        span.white--text {{ $t('common.pbehaviorPatterns') }}
      v-card
        v-card-text
          c-pbehavior-patterns-field(
            v-field="value.pbehavior_pattern",
            :required="isPatternRequired",
            :disabled="disabled",
            :name="entityPbehaviorFieldName",
            with-type,
            @input="errors.remove(entityPbehaviorFieldName)"
          )

    c-collapse-panel.mb-2(v-if="withEvent", :error="errors.has(eventPatternsFieldName)", color="grey")
      template(#header="")
        span.white--text {{ $t('common.eventPatterns') }}
      v-card
        v-card-text
          c-event-filter-patterns-field(
            v-field="value.event_pattern",
            :required="isPatternRequired",
            :disabled="disabled",
            :name="eventPatternsFieldName",
            @input="errors.remove(eventPatternsFieldName)"
          )
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
      const {
        alarm_pattern: alarmPattern,
        entity_pattern: entityPattern,
        pbehavior_pattern: pbehaviorPattern,
        event_pattern: eventPattern,
      } = this.value;

      return alarmPattern?.groups?.length
        || entityPattern?.groups?.length
        || pbehaviorPattern?.groups?.length
        || eventPattern?.groups?.length;
    },

    isPatternRequired() {
      return this.someRequired ? !this.hasPatterns : this.required;
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

    entityPbehaviorFieldName() {
      return this.preparePatternsFieldName('pbehavior_patterns');
    },
  },
  methods: {
    preparePatternsFieldName(name) {
      return [this.name, name].filter(Boolean).join('.');
    },
  },
};
</script>
