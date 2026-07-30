// validation.js
import { computed, watch } from 'vue';
import { debounce } from 'lodash';
import { useField } from 'vee-validate';

import { VALIDATION_DELAY } from '@/constants/common';

/**
 * Validation props.
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
  required: {
    type: Boolean,
    default: false,
  },
  rules: {
    type: [String, Object, Function],
    default: undefined,
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
  validationDelay: {
    type: Number,
    default: VALIDATION_DELAY,
  },
};

/**
 * Validation emits.
 */
export const validationEmits = ['update:modelValue', 'blur'];

/**
 * Base vee-validate registration hook.
 *
 * @param {Object} props
 * @param {Object} attrs
 * @returns {Object}
 */
export const useFieldValidation = (props, attrs = {}) => {
  const fieldName = computed(() => props.name || attrs.name || attrs.id || props.label || 'field');

  const fieldRules = computed(() => {
    if (props.rules) {
      return props.rules;
    }

    return {
      required: props.required || false,
    };
  });

  const field = useField(fieldName, fieldRules, {
    initialValue: props.modelValue,
    validateOnMount: props.validateOnMount,
    validateOnValueUpdate: false,
  });

  const validate = props.validationDelay
    ? debounce(field.validate, props.validationDelay)
    : field.validate;

  watch(
    () => props.modelValue,
    (value) => {
      field.setValue(value, false);

      if (props.validateOnValueUpdate) {
        validate();
      }
    },
  );

  return {
    ...field,
    validate,
    fieldName,
  };
};

/**
 * Vuetify field adapter over useFieldValidation.
 *
 * @param {Object} props
 * @param {Function} emit
 * @param {Object} attrs
 * @returns {Object}
 */
export const useVuetifyFieldValidation = (props, emit, attrs = {}) => {
  const field = useFieldValidation(props, attrs);

  const inputValue = computed({
    get: () => props.modelValue,
    set: (value) => {
      field.setValue(value, false);
      emit('update:modelValue', value);

      if (props.validateOnValueUpdate) {
        field.validate();
      }
    },
  });

  const errorMessages = computed(() => field.errors.value);

  const onBlur = (event, shouldValidate = props.validateOnBlur) => {
    field.handleBlur(event, shouldValidate);
    emit('blur', event);
  };

  return {
    ...field,
    inputValue,
    errorMessages,
    onBlur,
  };
};

/**
 * Array field validation hook.
 *
 * @param {Object} props
 * @param {Object} attrs
 * @returns {Object}
 */
export const useArrayFieldValidation = (props, attrs = {}) =>
  useFieldValidation(
    {
      ...props,
      rules:
        props.rules ||
        (value =>
          !props.required || (Array.isArray(value) && value.length > 0)
            ? true
            : unref(props.message) || false),
      validateOnValueUpdate: false,
    },
    attrs,
  );
