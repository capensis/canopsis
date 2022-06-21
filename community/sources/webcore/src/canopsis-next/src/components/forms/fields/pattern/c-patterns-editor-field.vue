<template lang="pug">
  v-layout(column)
    c-pattern-field.mb-2(
      v-if="withType",
      :value="patterns.id",
      :type="type",
      :disabled="disabled",
      return-object,
      required,
      @input="updatePattern"
    )

    v-flex
      v-alert.pre-wrap(v-if="errorMessage", value="true") {{ errorMessage }}

    v-tabs(
      v-if="!withType || patterns.id",
      v-model="activeTab",
      slider-color="primary",
      centered
    )
      v-tab(
        :disabled="!isSimpleTab && hasJsonError",
        :href="`#${$constants.PATTERN_EDITOR_TABS.simple}`"
      ) {{ $t('pattern.simpleEditor') }}
      v-tab-item(:value="$constants.PATTERN_EDITOR_TABS.simple")
        c-pattern-groups-field.mt-2(
          v-field="patterns.groups",
          :disabled="formDisabled",
          :name="patternGroupsFieldName",
          :type="type",
          :required="required",
          :attributes="attributes"
        )

      v-tab(:href="`#${$constants.PATTERN_EDITOR_TABS.advanced}`") {{ $t('pattern.advancedEditor') }}
      v-tab-item(:value="$constants.PATTERN_EDITOR_TABS.advanced", lazy)
        c-patterns-advanced-editor-field(
          :value="patternsJson",
          :disabled="disabled || !isCustomPattern",
          :attributes="attributes",
          :name="patternJsonFieldName",
          @input="updateGroupsFromPatterns"
        )

    v-layout(v-if="withType && !isCustomPattern", justify-end)
      v-btn.mx-0(
        color="primary",
        dark,
        @click="updatePatternToCustom"
      ) {{ $t('common.edit') }}
</template>

<script>
import { PATTERN_CUSTOM_ITEM_VALUE, PATTERN_EDITOR_TABS } from '@/constants';

import { formGroupsToPatternRules, patternsToGroups, patternToForm } from '@/helpers/forms/pattern';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      required: true,
    },
    attributes: {
      type: Array,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'patterns',
    },
    required: {
      type: Boolean,
      default: false,
    },
    type: {
      type: String,
      required: false,
    },
    withType: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      activeTab: PATTERN_EDITOR_TABS.simple,
      patternsJson: [],
    };
  },
  computed: {
    patternGroupsFieldName() {
      return `${this.name}.groups`;
    },

    patternJsonFieldName() {
      return `${this.name}.json`;
    },

    hasJsonError() {
      return this.errors.has(this.patternJsonFieldName);
    },

    errorMessage() {
      return this.errors.collect(this.name)?.join('\n');
    },

    isSimpleTab() {
      return this.activeTab === PATTERN_EDITOR_TABS.simple;
    },

    formDisabled() {
      return this.disabled || (this.withType && !this.isCustomPattern);
    },

    isCustomPattern() {
      return this.patterns.id === PATTERN_CUSTOM_ITEM_VALUE;
    },
  },
  watch: {
    activeTab(newTab) {
      if (newTab === PATTERN_EDITOR_TABS.advanced) {
        this.patternsJson = formGroupsToPatternRules(this.patterns.groups);
      }
    },
  },
  created() {
    this.$validator.attach({
      name: this.name,
      getter: () => this.patterns.length,
      context: () => this,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
  },
  methods: {
    updatePattern(pattern) {
      const { groups } = patternToForm(pattern);

      this.updateModel({
        ...this.patterns,
        id: pattern._id,
        groups,
      });
    },

    updatePatternToCustom() {
      this.updateField('id', PATTERN_CUSTOM_ITEM_VALUE);
    },

    updateGroupsFromPatterns(patterns) {
      this.updateField('groups', patternsToGroups(patterns));

      this.patternsJson = patterns;
    },
  },
};
</script>
