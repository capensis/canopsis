<template>
  <div class="pattern-rules-field">
    <v-row
      v-for="(rule, index) in rules"
      :key="rule.key"
      align="center"
      class="pattern-rules-field__row my-0"
      density="compact"
    >
      <v-col class="pattern-rules-field__rule py-0">
        <pattern-rule-field
          v-bind="getRuleProps(rule)"
          :rule="rule"
          :name="rule.key"
          :attributes="attributes"
          :disabled-field="rule.disabled?.field"
          :disabled-operator="rule.disabled?.operator"
          @update:rule="updateRule(index, $event)"
        />
      </v-col>
      <v-col class="pattern-rules-field__action py-0">
        <c-action-btn
          :tooltip="$t('pattern.removeRule')"
          :disabled="readonly || disabled || rule.disabled?.remove"
          type="delete"
          @click="removeItemFromArray(index)"
        />
      </v-col>
    </v-row>

    <c-btn-with-error
      v-if="!readonly"
      :error="hasRulesErrors ? $t('pattern.errors.existExcluded') : ''"
      :disabled="disabled"
      class="align-self-start"
      variant="outlined"
      @click="addFilterRule"
    >
      {{ $t('pattern.addRule') }}
    </c-btn-with-error>
  </div>
</template>

<script>
import { computed } from 'vue';

import {
  patternRuleToForm,
  convertValueByOperator,
  convertValueByType,
  getOperatorsByRule,
} from '@/helpers/entities/pattern/form';

import { useField } from '@/hooks/form';
import { formArrayMixin } from '@/mixins/form';

import PatternRuleField from './pattern-rule-field.vue';

export default {
  emits: ['update:rules'],
  components: { PatternRuleField },
  mixins: [formArrayMixin],
  props: {
    rules: {
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
      default: 'rules',
    },
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const fieldName = computed(() => props.name);
    const enabledRulesCount = () => props.rules.filter((rule) => {
        const attribute = props.attributes.find(({ value }) => value === rule.attribute);

        return !attribute?.options?.disabled;
      }).length;
    const { errorMessage, validate } = useField(
      fieldName,
      () => (enabledRulesCount() ? true : 'Required'),
      {
        initialValue: enabledRulesCount(),
        validateOnValueUpdate: false,
      },
    );

    return {
      rulesErrorMessage: errorMessage,
      validateRulesField: validate,
    };
  },
  computed: {
    rulesMap() {
      return this.attributes.reduce((acc, { value, options = {} }) => {
        acc[value] = options;

        return acc;
      }, {});
    },

    hasRulesErrors() {
      return !!this.rulesErrorMessage;
    },
  },
  watch: {
    rules() {
      this.$nextTick(this.validateRulesField);
    },
  },
  methods: {
    getOptionsByRule(rule) {
      return this.rulesMap[rule.attribute] ?? {};
    },

    getRuleProps(rule) {
      const { operators, type, disabled, ...props } = this.getOptionsByRule(rule);

      return {
        ...props,
        disabled: disabled || this.readonly || this.disabled,
        type,
        operators: operators ?? getOperatorsByRule(rule, type),
      };
    },

    getUpdatedRule(rule, newRule) {
      const { defaultValue, operators } = this.getRuleProps(newRule);

      const updatedRule = { ...newRule };

      if (updatedRule.attribute !== rule.attribute) {
        updatedRule.operator = '';
        updatedRule.field = '';
        updatedRule.dictionary = '';
        updatedRule.value = defaultValue;
      } else if (updatedRule.operator !== rule.operator) {
        updatedRule.value = convertValueByOperator(updatedRule.value, updatedRule.operator);
      }

      if (rule.fieldType !== updatedRule.fieldType) {
        updatedRule.value = convertValueByType(newRule.value, updatedRule.fieldType);
      }

      if (updatedRule.value !== rule.value && operators?.length === 1) {
        [updatedRule.operator] = operators;
      }

      if (!operators.includes(updatedRule.operator)) {
        updatedRule.operator = '';
      }

      return updatedRule;
    },

    updateRule(index, newRule) {
      const rule = this.rules[index];

      this.updateItemInArray(index, this.getUpdatedRule(rule, newRule));
    },

    addFilterRule() {
      const [firstAttribute] = this.attributes;

      this.addItemIntoArray(patternRuleToForm({ field: firstAttribute?.value }));
    },
  },
};
</script>

<style lang="scss" scoped>
.pattern-rules-field {
  display: flex;
  flex-direction: column;
  row-gap: 12px;
  min-width: 0;
  width: 100%;

  &__row {
    width: 100%;
  }

  &__rule {
    flex: 1 1 0;
    max-width: none;
  }

  &__action {
    display: flex;
    flex: 0 0 40px;
    justify-content: center;
    max-width: 40px;
    width: 40px;
  }
}
</style>
