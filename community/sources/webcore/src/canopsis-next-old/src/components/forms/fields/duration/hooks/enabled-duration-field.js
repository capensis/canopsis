import {
  computed,
  unref,
  inject,
  onMounted,
  onBeforeUnmount,
} from 'vue';

import { AVAILABLE_TIME_UNITS } from '@/constants';

import { convertUnit } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';
import { useComponentInstance } from '@/hooks/vue';

/**
 * Hook for handling enabled duration field functionality.
 *
 * @param {Object} options - The options object containing:
 * @param {Duration} options.duration - The duration object.
 * @param {DurationUnit[]} options.units - The units array.
 * @param {Duration} options.after - The after object.
 * @param {string} options.name - The name of the field.
 * @returns {Object} An object containing:
 * - `timeUnits`: A computed property returning an array of time units with translated text.
 * - `min`: A computed property calculating the minimum value based on duration and after values.
 * - `validate`: A function to validate the field.
 */
export const useEnabledDurationField = ({ duration, units: propsUnits, after, name }) => {
  const { tc } = useI18n();

  const timeUnits = computed(() => {
    const units = unref(propsUnits) || [
      AVAILABLE_TIME_UNITS.day,
      AVAILABLE_TIME_UNITS.week,
      AVAILABLE_TIME_UNITS.month,
      AVAILABLE_TIME_UNITS.year,
    ];

    return units.map(({ value, text }) => ({
      value,
      text: tc(text, unref(duration)?.value),
    }));
  });

  const min = computed(() => {
    const unwrappedDuration = unref(duration);
    const unwrappedAfter = unref(after);

    if (!unwrappedDuration?.enabled || !unwrappedAfter) {
      return 1;
    }

    return Math.floor(convertUnit(unwrappedAfter.value, unwrappedAfter.unit, unwrappedDuration.unit)) + 1;
  });

  const validator = inject('$validator');
  const instance = useComponentInstance();

  const attachField = () => {
    const fieldOptions = {
      name: unref(name),
      vm: instance,
      getter: () => unref(duration),
    };

    validator.attach(fieldOptions);
  };

  const detachField = () => validator.detach(unref(name));

  const validate = () => {
    const unwrappedName = unref(name);

    if (validator.errors.has(unwrappedName)) {
      return validator.validate(unwrappedName);
    }

    return true;
  };

  onMounted(() => attachField());
  onBeforeUnmount(() => detachField());

  return {
    timeUnits,
    min,
    validate,
  };
};
