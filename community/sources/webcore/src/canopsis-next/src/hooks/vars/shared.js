import { mapValues } from 'lodash';
import { unref } from 'vue';

import { varsToVariables } from '@/helpers/variables';

import { useI18n } from '@/hooks/i18n';

/**
 * Hook for preparing and translating template variables
 *
 * @param {string} prefix - Translation key prefix for variable text
 * @returns {Object} Object containing prepare function
 */
export const useVarsPrepare = (prefix) => {
  const { t } = useI18n();

  /**
   * Recursively adds translations to template variables
   *
   * @param {Array} items - Array of template variable items
   * @returns {Array} Array of items with translated text
   */
  const addTranslationToTemplateVars = (items = []) => items.map((item) => {
    const result = {
      ...item,

      text: item.alias ? item.text : t(`${prefix}.${item.text}`),
    };

    if (item.variables?.length) {
      result.variables = addTranslationToTemplateVars(item.variables);
    }

    return result;
  });

  /**
   * Prepares template variables by converting them to variables format and adding translations
   *
   * @param {Array} items - Array of template variable items to prepare
   * @returns {Object} Object with variable groups mapped to translated variables
   */
  const prepare = (items = []) => mapValues(varsToVariables(unref(items)), addTranslationToTemplateVars);

  return {
    prepare,
  };
};
