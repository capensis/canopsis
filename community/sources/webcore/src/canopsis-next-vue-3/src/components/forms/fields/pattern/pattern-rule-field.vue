<template>
  <div class="pattern-rule-field">
    <v-row class="pattern-rule-field__row" density="compact">
      <v-col :cols="firstColumnCols" class="pattern-rule-field__col">
        <v-row class="pattern-rule-field__row" density="compact">
          <v-col :cols="attributeColumnCols" class="pattern-rule-field__col">
            <pattern-attribute-field
              :value="rule.attribute"
              :items="attributes"
              :name="name"
              :disabled="disabled || disabledField"
              return-object
              required
              @update:value="updateAttribute"
            />
          </v-col>
          <v-col v-if="isAnyInfosRule" cols="8" class="pattern-rule-field__col">
            <c-infos-attribute-field
              :value="rule"
              :items="infos"
              :name="name"
              :disabled="disabled"
              :combobox="isInfosRule"
              row
              @update:value="updateModel"
            />
          </v-col>
          <v-col v-else-if="isObjectRule" cols="6" class="pattern-rule-field__col">
            <c-text-field
              :model-value="rule.dictionary"
              :name="objectDictionaryName"
              :disabled="disabled"
              :label="$t('common.dictionary')"
              required
              @update:model-value="updateField('dictionary', $event)"
            />
          </v-col>
        </v-row>
      </v-col>
      <v-col v-if="rule.attribute" :cols="secondColumnCols" class="pattern-rule-field__col">
        <v-row class="pattern-rule-field__row flex-nowrap" density="compact">
          <v-col v-if="isInfosValueField" class="pattern-rule-field__input-type-col">
            <c-input-type-field
              :value="rule.fieldType"
              :label="$t('common.type')"
              :types="inputTypesWithDefinedType"
              :disabled="disabled"
              :name="name"
              @update:value="updateField('fieldType', $event)"
            />
          </v-col>
          <v-col v-if="shownOperatorField" :cols="operatorCols" class="pattern-rule-field__col">
            <pattern-operator-field
              :value="rule.operator"
              :operators="operators"
              :disabled="disabled || disabledOperator"
              :name="operatorFieldName"
              required
              @update:value="updateField('operator', $event)"
            />
          </v-col>
          <v-col
            v-if="rule.operator && operatorHasValue"
            :cols="valueCols"
            :class="valueColumnClasses"
          >
            <component :is="valueComponent.is" v-bind="valueComponent.props" />
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    <v-row v-if="notDefinedType && !disabled" class="pattern-rule-field__row" density="compact">
      <v-col cols="12" class="pattern-rule-field__col">
        <c-alert type="warning">
          {{ $t('pattern.notDefinedType') }}
        </c-alert>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import { isFunction } from 'lodash';

import {
  PATTERN_FIELD_TYPES,
  PATTERN_RULE_INFOS_FIELDS,
  PATTERN_RULE_TYPES,
  DEFAULT_PATTERN_FIELD_TYPES,
} from '@/constants';

import {
  getFieldType,
  isDateRuleType,
  isDurationRuleType,
  isExtraInfosRuleType,
  isInfosRuleType,
  isObjectRuleType,
  isOperatorHasValue,
  isNumberRuleType,
} from '@/helpers/entities/pattern/form';

import { formMixin } from '@/mixins/form';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

import PatternAttributeField from './pattern-attribute-field.vue';
import PatternOperatorField from './pattern-operator-field.vue';
import PatternRuleFieldDateValue from './pattern-rule-field-date-value.vue';

export default {
  emits: ['update:rule'],
  components: {
    DateTimePickerTextField,
    PatternAttributeField,
    PatternOperatorField,
    PatternRuleFieldDateValue,
  },
  mixins: [formMixin],
  props: {
    rule: {
      type: Object,
      required: true,
    },
    attributes: {
      type: Array,
      default: () => [],
    },
    infos: {
      type: Array,
      default: () => [],
    },
    operators: {
      type: Array,
      default: () => [],
    },
    inputTypes: {
      type: Array,
      default: () => DEFAULT_PATTERN_FIELD_TYPES,
    },
    valueField: {
      type: Object,
      required: false,
    },
    type: {
      type: String,
      default: PATTERN_RULE_TYPES.string,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    disabledField: {
      type: Boolean,
      default: false,
    },
    disabledOperator: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'rule',
    },
    itemDisabled: {
      type: Function,
      required: false,
    },
  },
  computed: {
    operatorFieldName() {
      return `${this.name}.operator`;
    },

    valueFieldName() {
      return `${this.name}.value`;
    },

    inputType() {
      return getFieldType(this.rule.value);
    },

    isInfosRule() {
      return isInfosRuleType(this.type);
    },

    isExtraInfosRule() {
      return isExtraInfosRuleType(this.type);
    },

    isObjectRule() {
      return isObjectRuleType(this.type);
    },

    isAnyInfosRule() {
      return this.isInfosRule || this.isExtraInfosRule;
    },

    isAnyInfosRuleOrAlias() {
      return this.isAnyInfosRule || this.rule.alias;
    },

    isInfosValueField() {
      return this.rule.field === PATTERN_RULE_INFOS_FIELDS.value || this.rule.alias;
    },

    isDateRule() {
      return isDateRuleType(this.type) || this.rule.fieldType === PATTERN_FIELD_TYPES.timestamp;
    },

    isDurationRule() {
      return isDurationRuleType(this.type);
    },

    isNumberRule() {
      return isNumberRuleType(this.type) || this.rule.fieldType === PATTERN_FIELD_TYPES.number;
    },

    isBooleanRule() {
      return this.rule.fieldType === PATTERN_FIELD_TYPES.boolean;
    },

    inputTypesWithDefinedType() {
      const item = this.isInfosRule
        ? this.infos.find(info => info.value === this.rule.dictionary)
        : this.attributes.find(attribute => attribute.value === this.rule.attribute);

      if (!item?.definedType) {
        return this.inputTypes;
      }

      return this.inputTypes.map(type => ({
        ...type,
        defined: type.value === item.definedType,
      }));
    },

    valueComponent() {
      const valueProps = {
        value: this.rule.value,
        modelValue: this.rule.value,
        required: true,
        disabled: this.disabled,
        name: this.valueFieldName,
        label: this.$t('common.value'),
        'onUpdate:value': this.updateValue,
        'onUpdate:modelValue': this.updateValue,
      };

      if (this.isDateRule) {
        return {
          is: 'pattern-rule-field-date-value',
          props: {
            value: this.rule.range,
            operator: this.rule.operator,
            disabled: this.disabled,
            name: this.valueFieldName,
            'onUpdate:value': this.updateRange,
          },
        };
      }

      if (this.valueField) {
        const valueFieldProps = isFunction(this.valueField.props)
          ? this.valueField.props.call(this, this.rule)
          : this.valueField.props;

        const valueFieldOn = isFunction(this.valueField.on)
          ? this.valueField.on.call(this, this.rule)
          : this.valueField.on;

        return {
          is: this.valueField.is,
          props: {
            disabled: this.disabled,
            ...valueProps,
            ...valueFieldProps,
            ...valueFieldOn,
          },
        };
      }

      if (this.isDurationRule) {
        return {
          is: 'c-duration-field',
          props: {
            duration: this.rule.duration,
            disabled: this.disabled,
            name: this.valueFieldName,
            medium: true,
            'onUpdate:duration': this.updateDuration,
          },
        };
      }

      if (this.isNumberRule) {
        return {
          is: 'c-number-field',
          props: valueProps,
        };
      }

      if (this.isBooleanRule) {
        return {
          is: 'c-enabled-field',
          props: {
            ...valueProps,
            hideLabel: true,
          },
        };
      }

      return {
        is: 'c-mixed-input-field',
        props: {
          inputType: this.inputType,
          types: this.inputTypes,
          ...valueProps,
        },
      };
    },

    operatorHasValue() {
      const hasValue = isOperatorHasValue(this.rule.operator);

      if (this.isAnyInfosRule) {
        return this.isInfosValueField && hasValue;
      }

      return hasValue;
    },

    shownOperatorField() {
      return this.operators.length > 0 || this.operators[0] !== this.rule.operator;
    },

    objectDictionaryName() {
      return `${this.name}.dictionary`;
    },

    firstColumnCols() {
      return this.isAnyInfosRule ? 5 : 4;
    },

    secondColumnCols() {
      return this.isAnyInfosRuleOrAlias ? 7 : 8;
    },

    attributeColumnCols() {
      if (!this.isObjectRule && this.isAnyInfosRule) {
        return 4;
      }

      if (this.isObjectRule) {
        return 6;
      }

      return 12;
    },

    operatorCols() {
      if (this.isDateRule) {
        return 5;
      }

      if (this.isAnyInfosRule) {
        return 4;
      }

      return 6;
    },

    valueCols() {
      if (this.isAnyInfosRule || this.isDateRule) {
        return 7;
      }

      return 6;
    },

    valueColumnClasses() {
      return [
        'pattern-rule-field__col',
        {
          'pattern-rule-field__value-col--with-input-type':
            this.isAnyInfosRule && this.isInfosValueField,
        },
      ];
    },

    definedType() {
      return this.inputTypesWithDefinedType.find(type => type.defined);
    },

    notDefinedType() {
      return (
        this.isInfosValueField &&
        this.rule.fieldType &&
        this.definedType &&
        this.definedType.value !== this.rule.fieldType
      );
    },
  },
  methods: {
    preparedItemDisabled(item) {
      return this.itemDisabled?.(this.rule, item);
    },

    updateAttribute(attribute) {
      const newRule = {
        ...this.rule,

        attribute: attribute?.value,
        alias: attribute?.alias ?? false,
        fieldType: attribute?.definedType ?? PATTERN_FIELD_TYPES.string,
      };

      this.updateModel(newRule);
    },

    updateDuration(duration) {
      this.updateField('duration', duration);
    },

    updateRange(duration) {
      this.updateField('range', duration);
    },

    updateValue(value) {
      this.updateField('value', value);
    },
  },
};
</script>

<style lang="scss" scoped>
.pattern-rule-field {
  min-width: 0;
  width: 100%;

  &__col {
    flex-grow: 0;
    flex-shrink: 0;

    > * {
      width: 100%;
    }
  }

  &__input-type-col {
    flex: 0 0 40px;
    max-width: 40px;
    width: 40px;
  }

  &__value-col--with-input-type {
    flex: 1 1 0;
    max-width: none;
  }
}
</style>
