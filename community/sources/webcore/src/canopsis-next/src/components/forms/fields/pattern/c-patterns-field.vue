<template lang="pug">
  v-layout(column)
    c-collapse-panel.mb-2(v-if="alarm", color="grey")
      template(#header="")
        span.white--text {{ $t('common.alarmPatterns') }}
      v-card
        v-card-text
          c-alarm-patterns-field(
            v-field="value.alarm_pattern",
            :required="isPatternRequired",
            :disabled="disabled",
            :name="alarmPatternsFieldName",
            @input="errors.remove(alarmPatternsFieldName)"
          )

    c-collapse-panel.mb-2(v-if="entity", color="grey")
      template(#header="")
        span.white--text {{ $t('common.entityPatterns') }}
      v-card
        v-card-text
          c-entity-patterns-field(
            v-field="value.entity_pattern",
            :required="isPatternRequired",
            :disabled="disabled",
            :name="entityPatternsFieldName",
            @input="errors.remove(entityPatternsFieldName)"
          )

    c-collapse-panel(v-if="pbehavior", color="grey")
      template(#header="")
        span.white--text {{ $t('common.pbehaviorPatterns') }}
      v-card
        v-card-text
          c-pbehavior-patterns-field(
            v-field="value.pbehavior_pattern",
            :required="isPatternRequired",
            :disabled="disabled",
            :name="entityPbehaviorFieldName",
            @input="errors.remove(entityPbehaviorFieldName)"
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
    alarm: {
      type: Boolean,
      default: false,
    },
    entity: {
      type: Boolean,
      default: false,
    },
    pbehavior: {
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
      } = this.value;

      return alarmPattern.groups.length || entityPattern.groups.length || pbehaviorPattern.groups.length;
    },

    isPatternRequired() {
      return this.someRequired ? !this.hasPatterns : this.required;
    },

    alarmPatternsFieldName() {
      return this.preparePatternsFieldName('alarm_patterns');
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
