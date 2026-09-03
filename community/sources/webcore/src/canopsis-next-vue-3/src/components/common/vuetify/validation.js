import { debounce, isEqual, isNil, isNumber, isString } from 'lodash';
import { validate } from 'vee-validate';
import { computed, onBeforeUnmount, ref, watch } from 'vue';

import { useField, useFormContext } from '@/hooks/form';

import { VALIDATION_DELAY } from '@/constants/common';

/**
 * @typedef {boolean|string|Promise<boolean|string>} CustomValidationResult
 * Custom validator return value. `true`/`undefined` means valid; `false` or a string is an error.
 */

/**
 * @callback CustomValidationRule
 * Async or sync validator used in `validationRules`.
 *
 * @param {*} value Current field value.
 * @param {VuetifyFieldValidationProps} props Field component props.
 * @returns {CustomValidationResult}
 */

/**
 * @typedef {Object} VeeValidateRuleParams
 * Parameters passed to a vee-validate rule, optionally with a custom error message.
 *
 * @property {string} [message] Custom error message for the rule.
 */

/**
 * @typedef {Object<string, boolean|number|string|VeeValidateRuleParams>} VeeValidateRules
 * Map of vee-validate rule names to their parameters.
 */

/**
 * @typedef {Object} VeeValidateValidationContext
 * Context object provided by vee-validate to a field validator.
 *
 * @property {string} [name] Field name.
 * @property {string} [label] Field label used in generated messages.
 * @property {string} [field] Fallback field label when `label` is missing.
 */

/**
 * @callback FieldValidator
 * Validator function compatible with vee-validate `useField` rules.
 *
 * @param {*} value Current field value.
 * @param {VeeValidateValidationContext} ctx Validation context.
 * @returns {Promise<true|string>} `true` when valid, otherwise the first error message.
 */

/**
 * @typedef {VeeValidateRules|FieldValidator|FieldValidator[]} ValidationRulesPipeline
 * Rules passed to `useField`: a vee-validate rules object, a single validator, or a validator pipeline.
 */

/**
 * @typedef {Object} VuetifyFieldValidationProps
 * Resolved props consumed by `useVuetifyFieldValidation`.
 *
 * @property {*} [modelValue] Bound field value (`v-model`).
 * @property {string} [name] Field name used for native `name` and form registration.
 * @property {string} [label] Human-readable field label used in error messages.
 * @property {string} [density] Vuetify density (`compact` by default).
 * @property {boolean} [required] Whether the field is required.
 * @property {string} [requiredMessage] Custom message for the required rule. Supports `{field}`.
 * @property {number|string} [minLength] Minimum string length.
 * @property {string} [minLengthMessage] Custom min-length message. Supports `{field}` and `{min}`.
 * @property {number|string} [maxLength] Maximum string length.
 * @property {string} [maxLengthMessage] Custom max-length message. Supports `{field}` and `{max}`.
 * @property {number|string} [min] Minimum numeric value (fallback when `minValue` is empty).
 * @property {number|string} [minValue] Minimum numeric value (takes precedence over `min`).
 * @property {string} [minMessage] Custom message for `min`. Supports `{field}` and `{min}`.
 * @property {string} [minValueMessage] Custom message for `minValue`. Supports `{field}` and `{min}`.
 * @property {number|string} [max] Maximum numeric value (fallback when `maxValue` is empty).
 * @property {number|string} [maxValue] Maximum numeric value (takes precedence over `max`).
 * @property {string} [maxMessage] Custom message for `max`. Supports `{field}` and `{max}`.
 * @property {string} [maxValueMessage] Custom message for `maxValue`. Supports `{field}` and `{max}`.
 * @property {boolean} [decimal] Whether the value must be a decimal number.
 * @property {string} [decimalMessage] Custom decimal message. Supports `{field}`.
 * @property {boolean} [email] Whether the value must be a valid email.
 * @property {string} [emailMessage] Custom email message. Supports `{field}`.
 * @property {RegExp|string} [pattern] Regular expression used when `regex` is not set.
 * @property {RegExp|string} [regex] Regular expression (takes precedence over `pattern`).
 * @property {string} [patternMessage] Custom message for `pattern`. Supports `{field}`.
 * @property {string} [regexMessage] Custom message for `regex`. Supports `{field}`.
 * @property {boolean} [url] Whether the value must be a valid URL.
 * @property {string} [urlMessage] Custom URL message. Supports `{field}`.
 * @property {boolean|Function|Array} [unique] Uniqueness check: callback or list of existing values.
 * @property {string} [uniqueMessage] Custom unique message. Supports `{field}`.
 * @property {CustomValidationRule[]} [validationRules] Extra validators run after built-in rules.
 * @property {boolean} [validateOnMount] Whether to validate immediately on mount.
 * @property {boolean} [validateOnValueUpdate] Whether to validate after the value changes.
 * @property {boolean} [validateOnBlur] Whether to validate on blur.
 * @property {boolean} [bails] Whether to stop after the first failing rule.
 * @property {boolean} [standalone] Whether the field is independent from the parent form.
 * @property {number} [validationDelay] Debounce delay in milliseconds for value-update validation.
 */

/**
 * @typedef {Object} VuetifyFieldValidationReturn
 * Bindings returned by `useVuetifyFieldValidation` for Vuetify field components.
 *
 * @property {import('vue').ComputedRef<string|undefined>} inputName Native input `name` attribute.
 * @property {import('vue').WritableComputedRef<*>} inputValue Two-way field value bound to `v-model`.
 * @property {import('vue').ComputedRef<string>} errorMessage First internal validation error.
 * @property {import('vue').ComputedRef<string[]>} errorMessages Messages to pass to Vuetify (`error-messages`).
 * @property {Function} onBlur Blur handler that optionally validates and re-emits `blur`.
 */

/**
 * Converts a value to a number. Symbols become `NaN` instead of throwing.
 *
 * @param {*} value Value to coerce.
 * @returns {number}
 */
const toNumber = (value) => {
  if (typeof value === 'symbol') {
    return NaN;
  }

  return Number(value);
};

/**
 * Replaces `{key}` placeholders in a validation message with matching param values.
 *
 * @param {string} message Message template.
 * @param {Object<string, *>} params Placeholder values.
 * @returns {string}
 */
const formatMessage = (message, params) =>
  message.replace(/\{(\w+)}/g, (_, key) => params[key] ?? '');

/**
 * Resolves the display label used in generated validation messages.
 *
 * @param {VuetifyFieldValidationProps} props Field props.
 * @returns {string}
 */
const getFieldLabel = props => props.label || props.name || 'Field';

/**
 * Normalizes a custom validator result to vee-validate's `true | string` contract.
 *
 * @param {CustomValidationResult|undefined} result Raw validator result.
 * @returns {true|string}
 */
const normalizeValidationResult = (result) => {
  if (result === true || result === undefined) {
    return true;
  }

  if (result === false) {
    return 'Field is invalid';
  }

  return result;
};

/**
 * Checks whether a rule prop is present and not an empty string.
 *
 * @param {*} value Rule prop value.
 * @returns {boolean}
 */
const isRuleDefined = value => !isNil(value) && value !== '';

/**
 * Checks whether a numeric rule prop is present and finite.
 *
 * @param {*} value Rule prop value.
 * @returns {boolean}
 */
const isNumberRuleDefined = value => isRuleDefined(value) && Number.isFinite(toNumber(value));

/**
 * Checks whether a field value should be treated as empty.
 *
 * `null`, `undefined`, `''`, and empty arrays are empty. This keeps empty-to-empty
 * updates from being treated as user edits.
 *
 * @param {*} value Field value.
 * @returns {boolean}
 */
const isEmptyFieldValue = value =>
  isNil(value) || value === '' || (Array.isArray(value) && !value.length);

/**
 * Checks whether the field value actually changed for validation purposes.
 *
 * Empty-to-empty transitions are ignored so parent re-renders do not look like user input.
 *
 * @param {*} nextValue Incoming value.
 * @param {*} previousValue Current value.
 * @returns {boolean}
 */
const isFieldValueChanged = (nextValue, previousValue) =>
  !(isEmptyFieldValue(nextValue) && isEmptyFieldValue(previousValue)) &&
  !isEqual(nextValue, previousValue);

/**
 * Reads and interpolates a custom message prop, if one is set.
 *
 * @param {VuetifyFieldValidationProps} props Field props.
 * @param {string} messageProp Name of the message prop to read.
 * @param {Object<string, *>} [params={}] Extra placeholders besides `{field}`.
 * @returns {string|undefined}
 */
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

/**
 * Attaches a custom message to vee-validate rule parameters.
 *
 * Boolean rules become `{ message }`. Scalar rules become `{ [paramKey]: value, message }`.
 *
 * @param {boolean|number|string|Object} params Original rule parameters.
 * @param {string} [message] Custom error message.
 * @param {string} [paramKey] Parameter key used when `params` is a number or string.
 * @returns {boolean|number|string|VeeValidateRuleParams}
 */
const withMessage = (params, message, paramKey) => {
  if (!message) {
    return params;
  }

  if (params === true) {
    return { message };
  }

  if (paramKey && (isNumber(params) || isString(params))) {
    return message ? { [paramKey]: params, message } : { [paramKey]: params };
  }

  return { ...params, message };
};

/**
 * Builds the vee-validate rules object from field props.
 *
 * Only defined rule props are included. `minValue`/`maxValue` take precedence over `min`/`max`.
 * `regex` takes precedence over `pattern`.
 *
 * @param {VuetifyFieldValidationProps} props Field props.
 * @returns {VeeValidateRules}
 */
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
  if (isNumberRuleDefined(minNumericValue)) {
    const min = toNumber(minNumericValue);
    const messageProp = isRuleDefined(props.minValue) ? 'minValueMessage' : 'minMessage';

    rules.min_value = withMessage(min, getFieldMessage(props, messageProp, { min }), 'min');
  }

  const maxNumericValue = isRuleDefined(props.maxValue) ? props.maxValue : props.max;
  if (isNumberRuleDefined(maxNumericValue)) {
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

/**
 * Wraps built-in vee-validate rules in a single async validator.
 *
 * Used when custom `validationRules` are also present so both can run as a pipeline.
 *
 * @param {VeeValidateRules} rules Built-in vee-validate rules.
 * @param {VuetifyFieldValidationProps} props Field props.
 * @returns {FieldValidator}
 */
const createBuiltInRulesValidator = (rules, props) => async (value, ctx) => {
  const { valid, errors } = await validate(value, rules, {
    name: ctx.name,
    label: ctx.label ?? ctx.field,
    bails: props.bails,
  });

  return valid ? true : (errors[0] ?? true);
};

/**
 * Combines built-in props-based rules with custom `validationRules`.
 *
 * With no custom rules, returns a vee-validate rules object. Otherwise returns a validator
 * pipeline: built-in rules first, then each custom rule.
 *
 * @param {VuetifyFieldValidationProps} props Field props.
 * @returns {ValidationRulesPipeline}
 */
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

/**
 * Vue props that enable vee-validate on Vuetify field wrappers (`c-text-field`, `c-select`, …).
 *
 * Mix into a field component and pass the same props into `useVuetifyFieldValidation`.
 *
 * @type {Object}
 */
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
  density: {
    type: String,
    default: 'compact',
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

/**
 * Events emitted by fields that use `useVuetifyFieldValidation`.
 *
 * @type {string[]}
 */
export const validationEmits = ['update:modelValue', 'blur'];

/**
 * Wires a Vuetify field to vee-validate.
 *
 * Registers the field with the parent form, maps props to rules, debounces validation on value
 * updates, and hides internal errors until the user has edited the field, the form has been
 * submitted, or `validateOnMount` is set. External `error-messages` always take precedence.
 *
 * @param {VuetifyFieldValidationProps} props Field component props (`validationProps`).
 * @param {Function} emit Component emit function.
 * @param {Object} [attrs={}] Non-prop attributes; `name`, `id`, and `errorMessages` are read.
 * @returns {VuetifyFieldValidationReturn}
 */
export function useVuetifyFieldValidation(props, emit, attrs = {}) {
  const form = useFormContext();
  const valueChangedByUser = ref(false);

  /**
   * Error messages passed from the parent via `error-messages` / `errorMessages`.
   *
   * @type {import('vue').ComputedRef<string[]>}
   */
  const externalErrorMessages = computed(() => {
    const messages = attrs.errorMessages ?? attrs['error-messages'];

    if (!messages) {
      return [];
    }

    return Array.isArray(messages) ? messages : [messages];
  });

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

  /**
   * Runs field validation, debounced when `validationDelay` is greater than 0.
   *
   * @returns {Promise<Object>|undefined}
   */
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
      const hasChanged = isFieldValueChanged(nextValue, inputValue.value);

      if (hasChanged) {
        valueChangedByUser.value = true;
      }

      setValue(nextValue, false);
      emit('update:modelValue', nextValue);

      if (props.validateOnValueUpdate && valueChangedByUser.value) {
        internalValidate();
      }
    },
  });

  /**
   * Whether internal vee-validate errors should be shown.
   *
   * Errors stay hidden until the user edits the field, the form is submitted,
   * or `validateOnMount` is enabled.
   *
   * @type {import('vue').ComputedRef<boolean>}
   */
  const shouldShowInternalErrors = computed(() =>
    valueChangedByUser.value ||
    props.validateOnMount ||
    form.submitCount.value > 0,
  );

  const errorMessages = computed(() =>
    externalErrorMessages.value.length || !shouldShowInternalErrors.value
      ? externalErrorMessages.value
      : errors.value,
  );

  /**
   * Handles blur: optionally validates immediately and re-emits the event.
   *
   * @param {FocusEvent} event Native blur event.
   * @param {boolean} [shouldValidate=props.validateOnBlur] Whether to validate on this blur.
   */
  const onBlur = (event, shouldValidate = props.validateOnBlur) => {
    const shouldValidateField = shouldValidate && shouldShowInternalErrors.value;

    if (shouldValidateField) {
      debouncedValidate.cancel();
    }

    handleBlur(event, shouldValidateField);
    emit('blur', event);
  };

  watch(
    () => props.modelValue,
    (nextValue, prevValue) => {
      if (nextValue === prevValue) {
        return;
      }

      setValue(nextValue, false);

      if (props.validateOnValueUpdate && shouldShowInternalErrors.value) {
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
