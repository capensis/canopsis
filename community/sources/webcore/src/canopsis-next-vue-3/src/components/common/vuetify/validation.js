import { debounce, isNil, isNumber, isString } from 'lodash';
import { validate } from 'vee-validate';
import { computed, onBeforeUnmount, watch } from 'vue';

import { useField } from '@/hooks/form';

import { VALIDATION_DELAY } from '@/constants/common';

const toNumber = (value) => {
  if (typeof value === 'symbol') {
    return NaN;
  }

  return Number(value);
};

const formatMessage = (message, params) =>
  message.replace(/\{(\w+)}/g, (_, key) => params[key] ?? '');

const getFieldLabel = props => props.validationLabel || props.label || props.name || 'Field';

const normalizeValidationResult = (result) => {
  if (result === true || result === undefined) {
    return true;
  }

  if (result === false) {
    return 'Field is invalid';
  }

  return result;
};

const isRuleDefined = value => !isNil(value) && value !== '';

const getFieldMessage = (props, messageProp, params = {}) => {
  const message = props[messageProp];

  if (!message) {
    return undefined;
  }

  return formatMessage(message, {
    field: getFieldLabel(props),
    ...params,
  });
};

const withMessage = (params, message, paramKey) => {
  if (!message) {
    return params;
  }

  if (params === true) {
    return { message };
  }

  if (paramKey && (isNumber(params) || isString(params))) {
    return { [paramKey]: params, message };
  }

  return { ...params, message };
};

const buildFieldRules = (props) => {
  const rules = {};

  if (props.required) {
    rules.required = withMessage(true, getFieldMessage(props, 'requiredMessage'));
  }

  if (isRuleDefined(props.minLength)) {
    const length = toNumber(props.minLength);
    rules.min = withMessage(
      length,
      getFieldMessage(props, 'minLengthMessage', { min: length }),
      'length',
    );
  }

  if (isRuleDefined(props.maxLength)) {
    const length = toNumber(props.maxLength);
    rules.max = withMessage(
      length,
      getFieldMessage(props, 'maxLengthMessage', { max: length }),
      'length',
    );
  }

  const minNumericValue = isRuleDefined(props.minValue) ? props.minValue : props.min;
  if (isRuleDefined(minNumericValue)) {
    const min = toNumber(minNumericValue);
    const messageProp = isRuleDefined(props.minValue) ? 'minValueMessage' : 'minMessage';

    rules.min_value = withMessage(min, getFieldMessage(props, messageProp, { min }), 'min');
  }

  const maxNumericValue = isRuleDefined(props.maxValue) ? props.maxValue : props.max;
  if (isRuleDefined(maxNumericValue)) {
    const max = toNumber(maxNumericValue);
    const messageProp = isRuleDefined(props.maxValue) ? 'maxValueMessage' : 'maxMessage';

    rules.max_value = withMessage(max, getFieldMessage(props, messageProp, { max }), 'max');
  }

  if (props.decimal) {
    rules.decimal = withMessage(true, getFieldMessage(props, 'decimalMessage'));
  }

  if (props.email) {
    rules.email = withMessage(true, getFieldMessage(props, 'emailMessage'));
  }

  if (props.url) {
    rules.url = withMessage(true, getFieldMessage(props, 'urlMessage'));
  }

  const patternSource = props.regex ?? props.pattern;
  if (patternSource) {
    const messageProp = props.regex == null ? 'patternMessage' : 'regexMessage';

    rules.regex = withMessage({ regex: patternSource }, getFieldMessage(props, messageProp));
  }

  if (typeof props.unique === 'function') {
    rules.unique = withMessage({ validate: props.unique }, getFieldMessage(props, 'uniqueMessage'));
  } else if (Array.isArray(props.unique)) {
    rules.unique = withMessage({ values: props.unique }, getFieldMessage(props, 'uniqueMessage'));
  }

  return rules;
};

const createBuiltInRulesValidator = (rules, props) => async (value, ctx) => {
  const { valid, errors } = await validate(value, rules, {
    name: ctx.name,
    label: ctx.label ?? ctx.field,
    bails: props.bails,
  });

  return valid ? true : (errors[0] ?? true);
};

const createValidationRulesPipeline = (props) => {
  const builtInRules = buildFieldRules(props);
  const customRules = props.validationRules.map(
    rule => async value => normalizeValidationResult(await rule(value, props)),
  );

  if (!customRules.length) {
    return builtInRules;
  }

  const validators = [];

  if (Object.keys(builtInRules).length) {
    validators.push(createBuiltInRulesValidator(builtInRules, props));
  }

  validators.push(...customRules);

  return validators.length === 1 ? validators[0] : validators;
};

export const validationProps = {
  modelValue: {
    default: undefined,
  },
  name: {
    type: String,
    default: '',
  },
  label: {
    type: String,
    default: '',
  },
  validationLabel: {
    type: String,
    default: '',
  },
  required: {
    type: Boolean,
    default: false,
  },
  requiredMessage: {
    type: String,
    default: '',
  },
  minLength: {
    type: [Number, String],
    default: null,
  },
  minLengthMessage: {
    type: String,
    default: '',
  },
  maxLength: {
    type: [Number, String],
    default: null,
  },
  maxLengthMessage: {
    type: String,
    default: '',
  },
  min: {
    type: [Number, String],
    default: null,
  },
  minValue: {
    type: [Number, String],
    default: null,
  },
  minMessage: {
    type: String,
    default: '',
  },
  minValueMessage: {
    type: String,
    default: '',
  },
  max: {
    type: [Number, String],
    default: null,
  },
  maxValue: {
    type: [Number, String],
    default: null,
  },
  maxMessage: {
    type: String,
    default: '',
  },
  maxValueMessage: {
    type: String,
    default: '',
  },
  decimal: {
    type: Boolean,
    default: false,
  },
  decimalMessage: {
    type: String,
    default: '',
  },
  email: {
    type: Boolean,
    default: false,
  },
  emailMessage: {
    type: String,
    default: '',
  },
  pattern: {
    type: [RegExp, String],
    default: null,
  },
  regex: {
    type: [RegExp, String],
    default: null,
  },
  patternMessage: {
    type: String,
    default: '',
  },
  regexMessage: {
    type: String,
    default: '',
  },
  url: {
    type: Boolean,
    default: false,
  },
  urlMessage: {
    type: String,
    default: '',
  },
  unique: {
    type: [Boolean, Function, Array],
    default: false,
  },
  uniqueMessage: {
    type: String,
    default: '',
  },
  validationRules: {
    type: Array,
    default: () => [],
  },
  validateOnMount: {
    type: Boolean,
    default: false,
  },
  validateOnValueUpdate: {
    type: Boolean,
    default: true,
  },
  validateOnBlur: {
    type: Boolean,
    default: true,
  },
  bails: {
    type: Boolean,
    default: true,
  },
  standalone: {
    type: Boolean,
    default: false,
  },
  validationDelay: {
    type: Number,
    default: VALIDATION_DELAY,
  },
};

export const validationEmits = ['update:modelValue', 'blur'];

export function useVuetifyFieldValidation(props, emit, attrs = {}) {
  const inputName = computed(() => props.name || attrs.name || undefined);

  const fieldName = computed(() => props.name || attrs.name || attrs.id || props.label || 'field');

  const fieldLabel = computed(() => getFieldLabel(props));

  const fieldRules = computed(() => createValidationRulesPipeline(props));

  const {
    value,
    errors,
    errorMessage,
    handleBlur,
    setValue,
    validate: validateField,
  } = useField(fieldName, fieldRules, {
    initialValue: props.modelValue,
    label: fieldLabel,
    validateOnMount: props.validateOnMount,
    validateOnValueUpdate: false,
    bails: props.bails,
    standalone: props.standalone,
  });

  const debouncedValidate = debounce(validateField, props.validationDelay);

  const internalValidate = () => {
    if (!props.validationDelay) {
      debouncedValidate.cancel();

      return validateField();
    }

    return debouncedValidate();
  };

  onBeforeUnmount(() => {
    debouncedValidate.cancel();
  });

  const inputValue = computed({
    get: () => (props.modelValue !== undefined ? props.modelValue : value.value),
    set: (nextValue) => {
      setValue(nextValue, false);
      emit('update:modelValue', nextValue);

      if (props.validateOnValueUpdate) {
        internalValidate();
      }
    },
  });

  const errorMessages = computed(() => errors.value);

  const onBlur = (event, shouldValidate = props.validateOnBlur) => {
    if (shouldValidate) {
      debouncedValidate.cancel();
    }

    handleBlur(event, shouldValidate);
    emit('blur', event);
  };

  watch(
    () => props.modelValue,
    (nextValue, prevValue) => {
      if (nextValue === prevValue) {
        return;
      }

      setValue(nextValue, false);

      if (props.validateOnValueUpdate) {
        internalValidate();
      }
    },
  );

  watch(fieldRules, () => {
    if (errors.value.length) {
      internalValidate();
    }
  });

  return {
    inputName,
    inputValue,
    errorMessage,
    errorMessages,
    onBlur,
  };
}
