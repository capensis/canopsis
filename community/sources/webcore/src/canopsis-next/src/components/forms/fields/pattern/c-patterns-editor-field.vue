<template lang="pug">
  v-layout(column)
    c-pattern-field.mb-2(
      v-if="withType",
      :value="patterns.id",
      :type="type",
      return-object,
      required,
      @input="updatePattern"
    )

    v-tabs(v-if="!withType || patterns.id", slider-color="primary", centered)
      v-tab {{ $t('pattern.simpleEditor') }}
      v-tab-item
        c-pattern-groups-field.mt-2(
          v-field="patterns.groups",
          :disabled="formDisabled",
          :name="name",
          :type="type",
          :required="required",
          :attributes="attributes"
        )
      v-tab {{ $t('pattern.advancedEditor') }}
      v-tab-item(lazy)
        c-patterns-advanced-editor-field(
          :value="patternsJson",
          :disabled="disabled || !isCustomPattern",
          :attributes="attributes",
          name="advancedJson",
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
import { PATTERN_CUSTOM_ITEM_VALUE } from '@/constants';

import { formGroupsToPatternRules, patternsToGroups, patternToForm } from '@/helpers/forms/pattern';

import { formMixin } from '@/mixins/form';

export default {
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
      required: false,
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
  computed: {
    formDisabled() {
      return this.disabled || (this.withType && !this.isCustomPattern);
    },

    isCustomPattern() {
      return this.patterns.id === PATTERN_CUSTOM_ITEM_VALUE;
    },

    patternsJson() {
      return formGroupsToPatternRules(this.patterns.groups);
    },
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
    },
  },
};
</script>
