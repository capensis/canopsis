import { has, isPlainObject, isString } from 'lodash';
import { unref } from 'vue';

/**
 * Checks whether a server error value can be applied to a vee-validate field.
 *
 * @param {*} message Error value.
 * @returns {boolean}
 */
const isFieldErrorMessage = (message) => {
  if (isString(message)) {
    return true;
  }

  return Array.isArray(message) && message.every(isString);
};

/**
 * Removes vee-validate literal-path brackets from a field name.
 *
 * @param {string} [path=''] Field path.
 * @returns {string}
 */
export const unwrapFieldPath = (path = '') => String(path).replace(/^\[(.+)]$/i, '$1');

/**
 * Converts a dotted numeric path to vee-validate array syntax.
 * `tags.0.name` -> `tags[0].name`
 *
 * @param {string} [path=''] Field path.
 * @returns {string}
 */
export const toVeeFormPath = (path = '') => unwrapFieldPath(path).replace(/\.(\d+)(?=\.|$)/g, '[$1]');

/**
 * Collects registered vee-validate field paths from a form context.
 *
 * @param {Object} [formContext] Form context created by useForm.
 * @returns {Set<string>}
 */
const getRegisteredFieldPaths = (formContext) => {
  const registeredPaths = new Set();

  formContext
    ?.getAllPathStatesWithInstances?.()
    .forEach(state => registeredPaths.add(unref(state.path)));

  Object.keys(formContext?.fieldsInstances?.value ?? {}).forEach(path => registeredPaths.add(path));
  Object.keys(formContext?.errors?.value ?? {}).forEach(path => registeredPaths.add(path));

  return registeredPaths;
};

/**
 * Resolves a server field name to the path vee-validate actually registered.
 *
 * @param {Object} [formContext] Form context created by useForm.
 * @param {string} field Server field name.
 * @returns {string}
 */
export const resolveRegisteredFieldPath = (formContext, field) => {
  const unwrapped = unwrapFieldPath(field);
  const candidates = new Set([field, unwrapped, `[${unwrapped}]`, toVeeFormPath(field)]);
  const registeredPaths = getRegisteredFieldPaths(formContext);

  return (
    [...registeredPaths].find(
      path => candidates.has(path) || unwrapFieldPath(path) === unwrapped,
    ) ?? field
  );
};

/**
 * Normalizes a server error payload to a field-to-message map.
 *
 * Accepts both the raw API shape `{ errors: { directory: '...' } }` and the
 * already unwrapped interceptor payload `{ directory: '...' }`.
 *
 * @param {*} err Server error payload.
 * @returns {Object<string, string|string[]>}
 */
export const getServerFieldErrors = (err) => {
  if (!isPlainObject(err)) {
    return {};
  }

  if (isPlainObject(err.errors)) {
    return err.errors;
  }

  return err;
};

/**
 * Returns server errors that belong to the current form.
 *
 * @param {*} errors Server error payload.
 * @param {Object} [options]
 * @param {Object|import('vue').Ref<Object>} [options.formValues] Local form values.
 * @param {Object} [options.formContext] Form context created by useForm.
 * @param {Set<string>} [options.registeredPaths] Registered vee-validate field paths.
 * @returns {[string, string|string[]][]}
 */
export const getMatchedServerFieldErrors = (
  errors,
  { formValues, formContext, registeredPaths } = {},
) => {
  const fieldErrors = getServerFieldErrors(errors);
  const unwrappedFormValues = unref(formValues);
  const fieldPaths = registeredPaths ?? getRegisteredFieldPaths(formContext);
  const normalizedRegisteredPaths = new Set([...fieldPaths].map(unwrapFieldPath));

  return Object.entries(fieldErrors).filter(([field, message]) => {
    if (!isFieldErrorMessage(message)) {
      return false;
    }

    const unwrapped = unwrapFieldPath(field);

    return (
      normalizedRegisteredPaths.has(unwrapped) ||
      has(unwrappedFormValues, unwrapped) ||
      Array.isArray(unwrappedFormValues)
    );
  });
};

/**
 * Maps matching server field errors onto vee-validate and returns whether any were applied.
 *
 * @param {Object} formContext Form context created by useForm.
 * @param {*} [errors={}] Server error payload.
 * @param {string} [scope] Optional field name prefix.
 * @returns {boolean}
 */
export const applyServerFormErrors = (formContext, errors = {}, scope) => {
  if (!formContext?.addError) {
    return false;
  }

  const matchedErrors = getMatchedServerFieldErrors(errors, {
    formValues: formContext.values,
    formContext,
  });

  if (!matchedErrors.length) {
    return false;
  }

  matchedErrors.forEach(([field, message]) => {
    const path = scope ? `${scope}.${field}` : field;

    formContext.addError(path, message);
  });

  return true;
};

/**
 * Builds a popup message for a submit error that could not be mapped to fields.
 *
 * @param {*} err Server error payload.
 * @param {string} fallback Default message.
 * @returns {string}
 */
export const getSubmitErrorPopupText = (err, fallback) => {
  const fieldErrors = getServerFieldErrors(err);
  const message = Object.values(fieldErrors)
    .flatMap((value) => {
      if (isString(value)) {
        return [value];
      }

      if (Array.isArray(value)) {
        return value.filter(isString);
      }

      return [];
    })
    .join('\n');

  return message || err?.details || fallback;
};
