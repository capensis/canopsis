<template>
  <c-json-field
    :value="value"
    :label="$t('pattern.advancedEditor')"
    :readonly="disabled"
    :name="name"
    validate-on="button"
    rows="10"
    @update:value="updatePatternsFromJSON"
  >
    <template #append="{ jsonString }">
      <v-fade-transition>
        <pattern-fix-with-ai-btn v-if="errorMessage" :json-string="jsonString" />
      </v-fade-transition>
    </template>
  </c-json-field>
</template>

<script>
import { computed } from 'vue';
import { useField } from 'vee-validate';
import { isArray } from 'lodash';

import {
  isExtraInfosRuleType,
  isInfosRuleType,
  isValidPatternRule,
  isObjectRuleType,
} from '@/helpers/entities/pattern/form';

import { formMixin } from '@/mixins/form';

import PatternFixWithAiBtn from './pattern-fix-with-ai-btn.vue';

export default {
  emits: ['update:value'],
  components: {
    PatternFixWithAiBtn,
  },
  mixins: [formMixin],
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
  setup(props) {
    const fieldName = computed(() => props.name);
    const { errorMessage, setErrors } = useField(fieldName, undefined, {
      validateOnValueUpdate: false,
    });

    const setErrorMessage = message => setErrors(message ? [message] : []);

    return {
      errorMessage,
      setErrorMessage,
    };
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
          isInfosRuleType(options?.type) ||
          isExtraInfosRuleType(options?.type) ||
          isObjectRuleType(options?.type)
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
      return isArray(rules) && rules.length > 0 && rules.every(this.isValidPatternRule);
    },

    isValidPatterns(patterns) {
      return isArray(patterns) && patterns.every(this.isValidPatternRules);
    },

    updatePatternsFromJSON(patterns) {
      const isValidPatterns = this.isValidPatterns(patterns);

      if (isValidPatterns) {
        this.setErrorMessage(undefined);
        this.updateModel(patterns);
      } else {
        this.setErrorMessage(this.$t('pattern.errors.invalidPatterns'));
      }
    },
  },
};
</script>
