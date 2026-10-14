import { isNil } from 'lodash';
import { ref, unref, watch } from 'vue';

import { useModelField } from './model-field';

/**
 * Keeps an empty parent value selected in a lazy search field by mapping it to a sentinel.
 *
 * Clearing the field is kept locally until blur, then the empty default is restored
 * and emitted back to the parent. `0` is treated as a real value.
 *
 * @param {Object} options - Configuration options for the hook.
 * @param {Ref|*} options.value - Parent model value.
 * @param {symbol|string|number} [options.emptyValue] - Sentinel used by the field for the empty option.
 * @param {string|number|null} [options.emptyModelValue=''] - Value emitted to the parent for the empty option.
 * @param {Function} emit - Vue emit function.
 * @returns {{
 *   emptyValue: symbol|string|number,
 *   fieldValue: Ref,
 *   isEmptyValue: Function,
 *   toFieldValue: Function,
 *   toModelValue: Function,
 *   updateValue: Function,
 *   restoreEmptyDefault: Function,
 * }}
 */
export const useLazySearchEmptyDefaultValue = (
  {
    value,
    emptyValue = Symbol('empty-default'),
    emptyModelValue = '',
  },
  emit,
) => {
  const { updateModel } = useModelField({}, emit);

  /**
   * Checks whether a value represents the empty default option.
   *
   * @param {*} nextValue - Field or parent value.
   * @returns {boolean}
   */
  const isEmptyValue = nextValue =>
    isNil(nextValue) || nextValue === emptyModelValue || nextValue === emptyValue;

  /**
   * Maps an empty parent value to the sentinel used by the field.
   *
   * @param {*} nextValue - Selected field or parent value.
   * @returns {*}
   */
  const toFieldValue = nextValue => (isEmptyValue(nextValue) ? emptyValue : nextValue);

  /**
   * Maps the empty sentinel back to the value expected by the parent/API.
   *
   * @param {*} nextValue - Field value.
   * @returns {*}
   */
  const toModelValue = nextValue => (isEmptyValue(nextValue) ? emptyModelValue : nextValue);

  const fieldValue = ref(toFieldValue(unref(value)));

  /**
   * Updates the selected value. Empty values are kept locally until focus leaves the field.
   *
   * @param {*} nextValue - Selected field value.
   */
  const updateValue = (nextValue) => {
    if (isNil(nextValue)) {
      fieldValue.value = null;
      return;
    }

    fieldValue.value = nextValue;
    updateModel(toModelValue(nextValue));
  };

  /**
   * Restores the empty default when the field is cleared and then blurred / clicked outside.
   */
  const restoreEmptyDefault = () => {
    if (!isNil(fieldValue.value)) {
      return;
    }

    fieldValue.value = emptyValue;
    updateModel(emptyModelValue);
  };

  watch(() => unref(value), (nextValue) => {
    fieldValue.value = toFieldValue(nextValue);
  });

  return {
    emptyValue,
    fieldValue,
    isEmptyValue,
    toFieldValue,
    toModelValue,
    updateValue,
    restoreEmptyDefault,
  };
};
