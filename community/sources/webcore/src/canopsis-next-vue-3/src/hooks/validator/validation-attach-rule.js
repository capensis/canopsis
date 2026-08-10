import { useField } from 'vee-validate';
import { computed, ref, unref, watch } from 'vue';

import { useValidator } from './validator';

/**
 * Compatibility helper for old dynamic `validator.attach` calls.
 *
 * The field is registered synchronously through vee-validate, while `attachRule`
 * only enables already registered rules and supplies the legacy getter value.
 *
 * @param {string | import('vue').Ref<string>} name Field name.
 * @param {string | Object | import('vue').Ref<string | Object>} rules Vee-validate rules.
 * @param {Function} [defaultGetter] Field value getter.
 * @returns {Object}
 */
export const useValidationAttachRule = (name, rules, defaultGetter = () => undefined) => {
  const validator = useValidator();
  const attached = ref(false);
  const getter = ref(defaultGetter);

  const fieldName = computed(() => unref(name));
  const fieldRules = computed(() => (attached.value ? unref(rules) : undefined));

  const { validate, setValue } = useField(fieldName, fieldRules, {
    initialValue: defaultGetter(),
    validateOnValueUpdate: false,
  });

  /**
   * Updates the internal field value from the legacy getter.
   */
  const syncValue = () => setValue(getter.value(), false);

  watch(() => getter.value(), syncValue);

  /**
   * Enables the rule and stores the getter used before each validation.
   *
   * @param {Function} [nextGetter]
   */
  const attachRule = (nextGetter = defaultGetter) => {
    getter.value = nextGetter;
    attached.value = true;
    syncValue();
  };

  /**
   * Clears the current error and disables validation for the field.
   */
  const detachRule = () => {
    attached.value = false;
    validator.errors.clear(fieldName.value);
  };

  /**
   * Runs vee-validate validation with the latest getter value.
   *
   * @returns {Promise<boolean>}
   */
  const validateRule = async () => {
    syncValue();

    const { valid } = await validate();

    return valid;
  };

  return {
    validator,
    attachRule,
    detachRule,
    validateRule,
  };
};
