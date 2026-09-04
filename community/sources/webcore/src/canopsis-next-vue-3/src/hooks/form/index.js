import { get, isString, set } from 'lodash';
import {
  FormContextKey as VEE_FORM_CONTEXT_KEY,
  useField as useVeeField,
  useForm as useVeeForm,
  useFormErrors as veeUseFormErrors,
  useFieldError as veeUseFieldError,
} from 'vee-validate';
import { computed, inject, nextTick, onBeforeUnmount, provide, ref, unref, watch } from 'vue';

import { FORM_CONTEXT_KEY } from '@/constants';

import Observer from '@/services/observer';

import { scrollToFirstFormError } from '@/helpers/vuetify';

import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { usePendingHandler } from '@/hooks/query/pending';
import { useComponentInstance } from '@/hooks/vue';

import {
  applyServerFormErrors,
  getSubmitErrorPopupText,
  toLiteralVeeFieldPath,
  toVeeFormPath,
  unwrapFieldPath,
} from './server-errors';

export { useModelField } from './model-field';
export { useArrayModelField } from './array-model-field';
export { useLazySearch } from './lazy-search';
export {
  applyServerFormErrors,
  getServerFieldErrors,
  toLiteralVeeFieldPath,
} from './server-errors';

const fallbackFormContexts = new WeakMap();

/**
 * Creates form state separated from vee-validate.
 *
 * Values are stored only in a local ref and can be consumed by subforms through
 * model(path) and array(path). vee-validate is kept as the validation engine
 * for registered field wrappers and submit handling.
 *
 * @param {Object} [options={}] Vee-validate form options.
 * @param {Object} [options.initialValues={}] Initial form values.
 * @param {Function} [options.errorsToValidation=v => v] Maps a server error to field errors.
 * @returns {Object} Form context.
 */
export const useForm = ({
  initialValues = {},
  errorsToValidation = v => v,
  ...veeOptions
} = {}) => {
  const values = ref(initialValues);
  const fieldsInstances = ref({});

  const instance = useComponentInstance();
  const popups = usePopups();
  const { t } = useI18n();
  const afterSubmitObserver = new Observer();

  const veeForm = useVeeForm(veeOptions);
  const veeFormContext = instance?.provides?.[VEE_FORM_CONTEXT_KEY] ?? null;
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
    veeForm.getAllPathStates().map((state) => {
      const path = unref(state.path);

      return {
        ...state,
        instance: fieldsInstances.value[path] ?? fieldsInstances.value[unwrapFieldPath(path)],
      };
    });

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
   * Clears all field errors without resetting form values.
   */
  const resetErrors = () => {
    veeForm
      .getAllPathStates()
      .forEach(state => veeForm.setFieldError(unref(state.path), undefined));
    Object.keys(veeForm.errors.value).forEach(path => veeForm.setFieldError(path, undefined));
  };

  /**
   * Resets local form state and vee-validate validation state.
   *
   * Values stay in the local ref. vee-validate `resetForm({ values })` merges into its own
   * reactive object and can re-render children with missing nested models.
   *
   * @param {Object} [nextValues=initialValues] Next form values.
   */
  const reset = (nextValues = initialValues) => {
    values.value = nextValues;
    resetErrors();
    veeForm.setTouched(false);
    veeForm.submitCount.value = 0;
  };

  /**
   * Scrolls to the first invalid field visible in the current tab.
   */
  const scrollToFirstError = () => nextTick(() => scrollToFirstFormError(instance?.$el));

  /**
   * Sets a vee-validate error on a field and marks it as touched.
   * Accepts unwrapped server keys (`directory`) and resolves the registered path.
   *
   * @param {string} field Field path.
   * @param {string|string[]} message Error message.
   */
  const addError = (field, message) => {
    const unwrapped = unwrapFieldPath(field);
    const candidates = new Set([field, unwrapped, `[${unwrapped}]`, toVeeFormPath(field)]);
    const state = veeForm.getAllPathStates().find(item => candidates.has(unref(item.path)));
    const path = state ? unref(state.path) : toLiteralVeeFieldPath(field);

    veeForm.setFieldError(path, message);
    veeForm.setFieldTouched(path, true);
  };

  const toVeePath = (field) => {
    if (!isString(field)) {
      return field;
    }

    return toLiteralVeeFieldPath(field);
  };

  const setFieldError = (field, message) => veeForm.setFieldError(toVeePath(field), message);
  const validateField = (field, opts) => veeForm.validateField(toVeePath(field), opts);
  const setFieldTouched = (field, touched) => veeForm.setFieldTouched(toVeePath(field), touched);
  const resetField = (field, state) => veeForm.resetField(toVeePath(field), state);

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
    validateField,
    setFieldError,
    setErrors: veeForm.setErrors,
    setFieldTouched,
    resetField,
    resetErrors,
    addError,
    scrollToFirstError,
    errors: veeForm.errors,
    meta: veeForm.meta,
    submitCount: veeForm.submitCount,
    submitting,
    isSubmitting,
    veeFormContext,
  };

  /**
   * Maps server errors onto vee-validate fields and returns whether any were applied.
   *
   * @param {*} errors Server error payload.
   * @param {string} [scope] Optional field name prefix.
   * @returns {boolean}
   */
  const setFormErrors = (errors, scope) => applyServerFormErrors(form, errors, scope);

  /**
   * Creates a submit handler with shared form pending state and server error mapping.
   *
   * @param {Function} callback Valid submit callback.
   * @param {Function} [invalidCallback] Invalid submit callback.
   * @returns {Function} Submit handler.
   */
  const handleSubmit = (callback, invalidCallback) =>
    veeForm.handleSubmit(
      async (...args) => {
        try {
          const data = await runSubmit(() => callback(...args));

          await afterSubmitObserver.notify(data);

          return data;
        } catch (err) {
          const wasSet = setFormErrors(errorsToValidation(err));

          if (wasSet) {
            scrollToFirstError();

            return;
          }

          console.error(err);

          popups.error({
            text: getSubmitErrorPopupText(err, t('errors.default')),
          });
        }
      },
      (...args) => {
        scrollToFirstError();
        invalidCallback?.(...args);
      },
    );

  form.setFormErrors = setFormErrors;
  form.handleSubmit = handleSubmit;

  provide(FORM_CONTEXT_KEY, form);
  provide('$afterSubmitObserver', afterSubmitObserver);

  return form;
};

/**
 * Returns the nearest provided form context or creates a local one for the current component.
 *
 * @returns {Object} Form context created by useForm.
 */
export const useFormContext = () => {
  const form = inject(FORM_CONTEXT_KEY, null);

  if (form) {
    return form;
  }

  const instance = useComponentInstance();

  if (!instance) {
    return useForm();
  }

  if (!fallbackFormContexts.has(instance)) {
    fallbackFormContexts.set(instance, useForm());
  }

  return fallbackFormContexts.get(instance);
};

/**
 * Registers a vee-validate field and binds the current component instance to the form context.
 *
 * @param {string|import('vue').Ref<string>|import('vue').ComputedRef<string>} name Field name.
 * @param {*} [rules] Validation rules.
 * @param {Object} [options] Vee-validate field options.
 * @returns {Object} Vee-validate field API.
 */
export const useField = (name, rules, options = {}) => {
  const form = useFormContext();

  const fieldPath = computed(() => toLiteralVeeFieldPath(name));
  const field = useVeeField(fieldPath, rules, options);
  const instance = useComponentInstance();

  watch(
    () => unref(name),
    (fieldName, previousFieldName) => {
      if (previousFieldName && previousFieldName !== fieldName) {
        form.removeFieldInstance(previousFieldName);
      }

      form.addFieldInstance(fieldName, instance);
    },
    { immediate: true },
  );

  onBeforeUnmount(() => {
    form.removeFieldInstance(unref(name));
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
 * Returns current form errors with helpers compatible with old errors.has/collect usage.
 *
 * @returns {{
 *   errors: import('vue').ComputedRef<Object<string, string>>,
 *   getError: function(string|import('vue').Ref<string>|import('vue').ComputedRef<string>): string|undefined,
 *   hasError: function(string|import('vue').Ref<string>|import('vue').ComputedRef<string>): boolean,
 *   collectErrors: function(string|import('vue').Ref<string>|import('vue').ComputedRef<string>): string[],
 * }}
 */
export const useFormErrorsMap = () => {
  useFormContext();

  const errors = useFormErrors();
  const getErrorCandidates = (path) => {
    const field = unref(path);
    const unwrapped = unwrapFieldPath(field);

    return [
      field,
      toLiteralVeeFieldPath(field),
      unwrapped,
      `[${unwrapped}]`,
      toVeeFormPath(field),
    ].filter(Boolean);
  };
  const getError = (path) => {
    const candidates = getErrorCandidates(path);

    return candidates.map(candidate => errors.value[candidate]).find(Boolean);
  };
  const collectErrors = (path) => {
    const error = getError(path);

    return Array.isArray(error) ? error : [error].filter(Boolean);
  };
  const hasError = path => !!getError(path);

  return {
    errors,
    getError,
    hasError,
    collectErrors,
  };
};

/**
 * Returns the error message for a field path from the current form errors.
 *
 * @param {string|import('vue').Ref<string>|import('vue').ComputedRef<string>} path Field path.
 * @returns {import('vue').ComputedRef<string|undefined>} Field error message.
 */
export const useFieldError = (path) => {
  useFormContext();

  if (path == null) {
    return veeUseFieldError(path);
  }

  const fieldPath = computed(() => toLiteralVeeFieldPath(path));

  return veeUseFieldError(fieldPath);
};
