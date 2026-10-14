import { computed, unref } from 'vue';

import { GROUPED_TRIGGERS, PRO_TRIGGERS } from '@/constants';

import { isDeprecatedTrigger } from '@/helpers/entities/scenario/form';

import { useI18n } from '@/hooks/i18n';

export const DEFAULT_TRIGGERS_TRANSLATION_PREFIX = 'common.triggers';

/**
 * Prepare grouped trigger items for trigger select fields.
 *
 * @param {Object} options
 * @param {import('vue').Ref<string[]> | string[]} options.types
 * @param {import('vue').Ref<string> | string} options.translationKeyPrefix
 * @param {import('vue').Ref<boolean> | boolean} options.isProVersion
 * @param {import('vue').Ref<Object> | Object} options.preparedTriggersByTypes
 * @returns {{
 *   preparedTriggers: import('vue').ComputedRef<Object[]>,
 *   triggerToPreparedTrigger: (trigger?: Object) => Object,
 *   getTriggerSelectedText: (trigger?: Object, values?: Object) => string,
 * }}
 */
export const usePreparedTriggers = ({
  types,
  translationKeyPrefix,
  isProVersion,
  preparedTriggersByTypes,
}) => {
  const { t, te } = useI18n();

  const getTriggerTranslationKey = (type) => {
    const prefix = unref(translationKeyPrefix);
    const key = `${prefix}.${type}`;

    return prefix !== DEFAULT_TRIGGERS_TRANSLATION_PREFIX && te(`${key}.text`)
      ? key
      : `${DEFAULT_TRIGGERS_TRANSLATION_PREFIX}.${type}`;
  };

  /**
   * Add translated labels and deprecation state to a trigger.
   *
   * @param {Object} [trigger={}]
   * @returns {Object}
   */
  const triggerToPreparedTrigger = (trigger = {}) => {
    const { type } = trigger;
    const triggerKey = getTriggerTranslationKey(type);

    return {
      ...trigger,
      text: te(`${triggerKey}.text`) ? t(`${triggerKey}.text`) : type,
      helpText: te(`${triggerKey}.helpText`) ? t(`${triggerKey}.helpText`) : '',
      deprecated: isDeprecatedTrigger(type),
    };
  };

  /**
   * Get the label used for a selected trigger chip.
   *
   * @param {Object} [trigger={}]
   * @param {Object} [values={}]
   * @returns {string}
   */
  const getTriggerSelectedText = (trigger = {}, values = {}) => {
    const messageKey = `${getTriggerTranslationKey(trigger.type)}.selectedText`;

    return te(messageKey) ? t(messageKey, values) : trigger.text;
  };

  /**
   * Build a prepared trigger item, or skip a Pro-only type when Pro is unavailable.
   *
   * @param {string} type
   * @returns {Object|null}
   */
  const prepareTriggerItem = (type) => {
    if (PRO_TRIGGERS.includes(type) && !unref(isProVersion)) {
      return null;
    }

    return triggerToPreparedTrigger({
      ...unref(preparedTriggersByTypes)[type],
      type,
    });
  };

  const preparedTriggers = computed(() => {
    const typesList = unref(types);

    if (unref(translationKeyPrefix) !== DEFAULT_TRIGGERS_TRANSLATION_PREFIX) {
      return typesList.map(prepareTriggerItem).filter(Boolean);
    }

    const typesSet = new Set(typesList);

    return Object.entries(GROUPED_TRIGGERS).reduce((acc, [group, groupTypes]) => {
      const preparedItems = groupTypes
        .filter(type => typesSet.has(type))
        .map(prepareTriggerItem)
        .filter(Boolean);

      if (preparedItems.length) {
        const title = t(`common.triggersGroups.${group}`);

        acc.push({ header: title }, ...preparedItems);
      }

      return acc;
    }, []);
  });

  return {
    preparedTriggers,
    triggerToPreparedTrigger,
    getTriggerSelectedText,
  };
};
