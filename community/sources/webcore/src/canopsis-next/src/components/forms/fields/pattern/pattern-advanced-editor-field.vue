<template>
  <c-json-field
    :value="value"
    :label="$t('pattern.advancedEditor')"
    :readonly="disabled"
    :name="name"
    validate-on="button"
    rows="10"
    @input="updatePatternsFromJSON"
  />
</template>

<script>
import { isArray } from 'lodash';

import {
  isExtraInfosRuleType,
  isInfosRuleType,
  isValidPatternRule,
  isObjectRuleType,
} from '@/helpers/entities/pattern/form';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
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
  },
  computed: {
    disabledAttributes() {
      return this.attributes.reduce((acc, { value, options }) => {
        if (options?.disabled) {
          acc.push(value);
        }

        return acc;
      }, []);
    },
  },
  methods: {
    isDisabledField(field) {
      return this.disabledAttributes.includes(field);
    },

    isValidRuleField({ field }) {
      return this.attributes.some(({ value, options }) => {
        if (this.isDisabledField(value)) {
          return false;
        }

        if (
          isInfosRuleType(options?.type)
          || isExtraInfosRuleType(options?.type)
          || isObjectRuleType(options?.type)
        ) {
          return field.startsWith(`${value}.`);
        }

        return value === field;
      });
    },

    isValidPatternRule(rule) {
      return isValidPatternRule(rule) && this.isValidRuleField(rule);
    },

    isValidPatternRules(rules) {
      return isArray(rules)
        && rules.length > 0
        && rules.every(this.isValidPatternRule);
    },

    isValidPatterns(patterns) {
      return isArray(patterns) && patterns.every(this.isValidPatternRules);
    },

    updatePatternsFromJSON(patterns) {
      const isValidPatterns = this.isValidPatterns(patterns);

      if (isValidPatterns) {
        this.updateModel(patterns);
      } else {
        /**
         * Validator.errors.reset using nextTick twice, because we need to add errors after reset
         */
        setTimeout(() => {
          this.errors.add({ field: this.name, msg: this.$t('pattern.errors.invalidPatterns') });
        });
      }
    },
  },
};
</script>
