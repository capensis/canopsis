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
        c-json-field(
          :value="patternsJson",
          :label="$t('pattern.advancedEditor')",
          :readonly="disabled || !isCustomPattern",
          name="advancedJson",
          validate-on="button",
          rows="10",
          @input="updatePatternFromJSON"
        )

    v-layout(v-if="withType && !isCustomPattern", justify-end)
      v-btn.mx-0(
        color="primary",
        dark,
        @click="updatePatternToCustom"
      ) {{ $t('common.edit') }}
</template>

<script>
import { isNull, isNumber, isObject, isString } from 'lodash';

import { PATTERN_CONDITIONS, PATTERN_CUSTOM_ITEM_VALUE } from '@/constants';

import { formGroupsToPatternRules, isDatePatternRule, patternsToGroups, patternToForm } from '@/helpers/forms/pattern';

import { formMixin } from '@/mixins/form';
import {
  getFieldType,
  isExtraInfosRuleType,
  isInfosRuleType,
  isStringArrayFieldType,
} from '@/helpers/pattern';

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

    isValidRuleField({ field }) {
      return this.attributes.some(({ value, options }) => {
        if (isInfosRuleType(options?.type) || isExtraInfosRuleType(options?.type)) {
          return field.startsWith(value);
        }

        return value === field;
      });
    },

    isValidRuleConditionType({ cond }) {
      return Object.values(PATTERN_CONDITIONS).includes(cond?.type);
    },

    isValidRuleFieldType({ cond, field_type: fieldType }) {
      return !fieldType || getFieldType(cond?.value) === fieldType;
    },

    isValidRuleConditionValue({ field, cond, field_type: fieldType }) {
      if (isStringArrayFieldType(fieldType)) {
        return cond.value.every(isString);
      }

      if (!fieldType) {
        if (isObject(cond.value)) {
          switch (cond.type) {
            case PATTERN_CONDITIONS.absoluteTime:
              return isDatePatternRule(field) && isNumber(cond.value.from) && isNumber(cond.value.to);
            case PATTERN_CONDITIONS.greater:
            case PATTERN_CONDITIONS.less:
              return isNumber(cond.value.value) && isString(cond.value.unit);
          }

          return false;
        }
      }

      return !isNull(cond?.value);
    },

    updatePatternFromJSON(patterns) {
      const isValidGroups = patterns.every(rules => rules.every((rule) => {
        const isValidField = this.isValidRuleField(rule);
        const isValidFieldType = this.isValidRuleFieldType(rule);
        const isValidConditionType = this.isValidRuleConditionType(rule);
        const isValidConditionValue = this.isValidRuleConditionValue(rule);

        return isValidField && isValidFieldType && isValidConditionType && isValidConditionValue;
      }));

      if (isValidGroups) {
        this.updateField('groups', patternsToGroups(patterns));
      }
    },
  },
};
</script>
