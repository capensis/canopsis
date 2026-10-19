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
          :key="`${rule.key}-${rule.attribute}`"
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
import { computed, nextTick, toRef, watch } from 'vue';

import {
  patternRuleToForm,
  convertValueByOperator,
  convertValueByType,
  getOperatorsByRule,
} from '@/helpers/entities/pattern/form';

import { useField } from '@/hooks/form';
import { useArrayModelField } from '@/hooks/form/array-model-field';

import PatternRuleField from './pattern-rule-field.vue';

export default {
  emits: ['update:rules'],
  components: { PatternRuleField },
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
  setup(props, { emit }) {
    const { addItemIntoArray, updateItemInArray, removeItemFromArray } = useArrayModelField(
      props,
      emit,
    );
    const fieldName = toRef(props, 'name');
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

    const rulesMap = computed(() =>
      props.attributes.reduce((acc, { value, options = {} }) => {
        acc[value] = options;

        return acc;
      }, {}));

    const hasRulesErrors = computed(() => !!errorMessage.value);

    const getOptionsByRule = rule => rulesMap.value[rule.attribute] ?? {};

    const getRuleProps = (rule) => {
      const {
        operators,
        type,
        disabled,
        ...ruleProps
      } = getOptionsByRule(rule);

      return {
        ...ruleProps,
        disabled: disabled || props.readonly || props.disabled,
        type,
        operators: operators ?? getOperatorsByRule(rule, type),
      };
    };

    const getUpdatedRule = (rule, newRule) => {
      const { defaultValue, operators } = getRuleProps(newRule);

      const updatedRule = { ...newRule };

      if (updatedRule.attribute !== rule.attribute) {
        updatedRule.operator = '';
        updatedRule.field = '';
        updatedRule.dictionary = '';
        updatedRule.value = defaultValue;
      } else if (updatedRule.operator !== rule.operator) {
        updatedRule.value = convertValueByOperator(
          updatedRule.value,
          updatedRule.operator,
          updatedRule,
          rule,
        );
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
    };

    const updateRule = (index, newRule) => {
      const rule = props.rules[index];

      updateItemInArray(index, getUpdatedRule(rule, newRule));
    };

    const addFilterRule = () => {
      const [firstAttribute] = props.attributes;

      addItemIntoArray(patternRuleToForm({ field: firstAttribute?.value }));
    };

    watch(
      () => props.rules,
      () => nextTick(validate),
    );

    return {
      hasRulesErrors,
      getRuleProps,
      updateRule,
      addFilterRule,
      removeItemFromArray,
    };
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
