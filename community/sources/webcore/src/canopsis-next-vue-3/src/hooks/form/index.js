import { get, set } from 'lodash';
import {
  useField as useVeeField,
  useForm as useVeeForm,
  useFormErrors as veeUseFormErrors,
  useFieldError as veeUseFieldError,
} from 'vee-validate';
import { computed, inject, onBeforeUnmount, provide, ref, unref, watch } from 'vue';

import { usePendingHandler } from '@/hooks/query/pending';
import { useComponentInstance } from '@/hooks/vue';

export { useModelField } from './model-field';
export { useArrayModelField } from './array-model-field';
export { useLazySearch } from './lazy-search';

export const FORM_CONTEXT_KEY = Symbol('formContext');

/**
 * Creates form state separated from vee-validate.
 *
 * Values are stored only in a local ref and can be consumed by subforms through
 * model(path) and array(path). vee-validate is kept as the validation engine
 * for registered field wrappers and submit handling.
 *
 * @param {Object} [options={}] Vee-validate form options.
 * @param {Object} [options.initialValues={}] Initial form values.
 * @returns {Object} Form context.
 */
export const useForm = ({ initialValues = {}, ...veeOptions } = {}) => {
  const values = ref(initialValues);
  const fieldsInstances = ref({});

  const veeForm = useVeeForm(veeOptions);
  const { pending: submitting, handler: runSubmit } = usePendingHandler(callback => callback());
  const isSubmitting = computed(() => veeForm.isSubmitting.value || submitting.value);

  /**
   * Registers a Vue component instance for a field path.
   *
   * @param {string} name Field path.
   * @param {Vue} instance Field component instance.
   */
  const addFieldInstance = (name, instance) => {
    fieldsInstances.value[name] = instance;
  };

  /**
   * Unregisters a Vue component instance for a field path.
   *
   * @param {string} name Field path.
   */
  const removeFieldInstance = (name) => {
    delete fieldsInstances.value[name];
  };

  /**
   * Returns vee-validate path states with their registered component instances.
   *
   * @returns {Array<Object>} Path states with an `instance` property.
   */
  const getAllPathStatesWithInstances = () =>
    veeForm.getAllPathStates().map(state => ({
      ...state,
      instance: fieldsInstances.value[state.path],
    }));

  /**
   * Returns a nested value from the local form state.
   *
   * @param {string|import('vue').Ref<string>} path Value path.
   * @returns {*} Value by path.
   */
  const getValue = path => get(values.value, unref(path));

  /**
   * Updates a nested value in local form state.
   *
   * @param {string|import('vue').Ref<string>} path Value path.
   * @param {*} value Next value.
   */
  const setValue = (path, value) => {
    const resolvedPath = unref(path);

    set(values.value, resolvedPath, value);
  };

  /**
   * Creates a writable computed model for a nested form value.
   *
   * @param {string|import('vue').Ref<string>} path Value path.
   * @returns {import('vue').WritableComputedRef<*>} Writable model.
   */
  const model = path =>
    computed({
      get: () => getValue(path),
      set: value => setValue(path, value),
    });

  /**
   * Creates multiple writable computed models from field paths.
   *
   * Supports two forms:
   * - array form: keys are equal to paths;
   * - object form: keys are custom aliases, values are paths.
   *
   * @param {string[]|Object<string, string|import('vue').Ref<string>|import('vue').ComputedRef<string>>} fields
   * Field paths or aliases mapped to field paths.
   *
   * @returns {Object<string, import('vue').WritableComputedRef<*>>}
   * Object where every key contains a writable computed model.
   *
   * @example
   * // Simple fields
   * const {
   *   title,
   *   email,
   * } = form.models(['title', 'email']);
   *
   * title.value = 'New title';
   * email.value = 'test@example.com';
   *
   * @example
   * // Nested fields with aliases
   * const {
   *   userEmail,
   *   userName,
   * } = form.models({
   *   userEmail: 'user.email',
   *   userName: 'user.name',
   * });
   *
   * userEmail.value = 'test@example.com';
   * userName.value = 'John';
   *
   * @example
   * // Usage in template
   * const {
   *   title,
   *   description,
   * } = form.models(['title', 'description']);
   *
   * // <c-text-field v-model="title" name="title" />
   * // <c-text-field v-model="description" name="description" />
   */
  const models = (fields) => {
    if (Array.isArray(fields)) {
      return fields.reduce((acc, path) => {
        acc[path] = model(path);

        return acc;
      }, {});
    }

    return Object.entries(fields).reduce((acc, [key, path]) => {
      acc[key] = model(path);

      return acc;
    }, {});
  };

  /**
   * Creates a writable computed model for a mutable array value.
   *
   * @param {string|import('vue').Ref<string>} path Array path.
   * @returns {import('vue').WritableComputedRef<Array>} Writable array model.
   */
  const array = path =>
    computed({
      get: () => {
        const current = getValue(path);

        if (Array.isArray(current)) {
          return current;
        }

        const nextValue = [];

        setValue(path, nextValue);

        return nextValue;
      },
      set: value => setValue(path, Array.isArray(value) ? value : []),
    });

  /**
   * Appends an item to an array value.
   *
   * @param {string|import('vue').Ref<string>} path Array path.
   * @param {*} item Item to append.
   */
  const push = (path, item) => {
    array(path).value.push(item);
  };

  /**
   * Removes an item from an array value by index.
   *
   * @param {string|import('vue').Ref<string>} path Array path.
   * @param {number} index Item index.
   */
  const remove = (path, index) => {
    array(path).value.splice(index, 1);
  };

  /**
   * Replaces an item in an array value by index.
   *
   * @param {string|import('vue').Ref<string>} path Array path.
   * @param {number} index Item index.
   * @param {*} value Next item value.
   */
  const replace = (path, index, value) => {
    array(path).value[index] = value;
  };

  /**
   * Updates a nested field inside an array item.
   *
   * @param {string|import('vue').Ref<string>} path Array path.
   * @param {number} index Item index.
   * @param {string} field Nested field path inside item.
   * @param {*} value Next field value.
   */
  const update = (path, index, field, value) => {
    const current = array(path).value;

    if (!current[index]) {
      return;
    }

    set(current[index], field, value);
  };

  /**
   * Resets local form state and vee-validate validation state.
   *
   * @param {Object} [nextValues=initialValues] Next form values.
   */
  const reset = (nextValues = initialValues) => {
    values.value = nextValues;

    veeForm.resetForm();
  };

  /**
   * Creates a submit handler with shared form pending state.
   *
   * @param {Function} callback Valid submit callback.
   * @param {Function} [invalidCallback] Invalid submit callback.
   * @returns {Function} Submit handler.
   */
  const handleSubmit = (callback, invalidCallback) =>
    veeForm.handleSubmit((...args) => runSubmit(() => callback(...args)), invalidCallback);

  /**
   * Clears all field errors without resetting form values.
   */
  const resetErrors = () => {
    veeForm
      .getAllPathStates()
      .forEach(state => veeForm.setFieldError(unref(state.path), undefined));
    Object.keys(veeForm.errors.value).forEach(path => veeForm.setFieldError(path, undefined));
  };

  const form = {
    values,

    model,
    models,
    array,
    getValue,
    setValue,

    push,
    remove,
    replace,
    update,
    reset,

    fieldsInstances,
    addFieldInstance,
    removeFieldInstance,
    getAllPathStatesWithInstances,

    validate: veeForm.validate,
    validateField: veeForm.validateField,
    setFieldError: veeForm.setFieldError,
    setErrors: veeForm.setErrors,
    setFieldTouched: veeForm.setFieldTouched,
    resetField: veeForm.resetField,
    resetErrors,
    handleSubmit,
    errors: veeForm.errors,
    meta: veeForm.meta,
    submitCount: veeForm.submitCount,
    submitting,
    isSubmitting,
  };

  provide(FORM_CONTEXT_KEY, form);

  return form;
};

/**
 * Returns the nearest provided form context.
 *
 * @returns {Object} Form context created by useForm.
 * @throws {Error} When called outside of a form context.
 */
export const useFormContext = () => {
  const form = inject(FORM_CONTEXT_KEY);

  if (!form) {
    throw new Error('Missing form context');
  }

  return form;
};

/**
 * Registers a vee-validate field and binds the current component instance to the form context.
 *
 * @param {string|import('vue').Ref<string>|import('vue').ComputedRef<string>} name Field name.
 * @param {*} [rules] Validation rules.
 * @param {Object} [options] Vee-validate field options.
 * @returns {Object} Vee-validate field API.
 */
export const useField = (name, rules, options) => {
  const field = useVeeField(name, rules, options);
  const form = inject(FORM_CONTEXT_KEY, null);
  const instance = useComponentInstance();

  watch(
    () => unref(name),
    (fieldName, previousFieldName) => {
      if (!form) {
        return;
      }

      if (previousFieldName && previousFieldName !== fieldName) {
        form.removeFieldInstance(previousFieldName);
      }

      form.addFieldInstance(fieldName, instance);
    },
    { immediate: true },
  );

  onBeforeUnmount(() => {
    if (form) {
      form.removeFieldInstance(unref(name));
    }
  });

  return field;
};

/**
 * Shortcut for creating a writable model from the current form context.
 *
 * @param {string|import('vue').Ref<string>} path Value path.
 * @returns {import('vue').WritableComputedRef<*>} Writable model.
 */
export const useFormModel = path => useFormContext().model(path);

/**
 * Creates an object of writable computed models from the current form context.
 *
 * @param {string[]|Object<string, string|import('vue').Ref<string>|import('vue').ComputedRef<string>>} fields
 * @returns {Object<string, import('vue').WritableComputedRef<*>>}
 */
export const useFormModels = (fields = []) => {
  const form = useFormContext();

  return form.models(fields);
};

/**
 * Shortcut for creating a writable array model from the current form context.
 *
 * @param {string|import('vue').Ref<string>} path Array path.
 * @returns {import('vue').WritableComputedRef<Array>} Writable array model.
 */
export const useFormArray = path => useFormContext().array(path);

export const useFormErrors = veeUseFormErrors;

/**
 * Returns the error message for a field path from the current form errors.
 *
 * @param {string|import('vue').Ref<string>|import('vue').ComputedRef<string>} path Field path.
 * @returns {import('vue').ComputedRef<string|undefined>} Field error message.
 */
export const useFieldError = veeUseFieldError;
