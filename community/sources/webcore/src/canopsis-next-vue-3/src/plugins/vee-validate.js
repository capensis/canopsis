import { isEmpty } from 'lodash';
import { watch, unref } from 'vue';
import { configure, defineRule } from 'vee-validate';
import { all as veeValidateRules } from '@vee-validate/rules';
import { localize, setFallbackLocale, setLocale } from '@vee-validate/i18n';

import { DEFAULT_LOCALE } from '@/config';

import { veeValidateLocales } from './vee-validate/locales';

/**
 * Returns a per-field custom validation message when it is passed through rule params.
 *
 * @param {import('vee-validate').FieldValidationMetaInfo} ctx
 * @returns {string | undefined}
 */
const getCustomRuleMessage = (ctx) => {
  const params = ctx.rule?.params;

  if (params && !Array.isArray(params) && params.message) {
    return params.message;
  }

  return undefined;
};

/**
 * Creates vee-validate message generator backed by `@vee-validate/i18n`
 * and keeps it in sync with the application vue-i18n locale.
 *
 * Custom field messages from wrapper props take priority over localized defaults.
 *
 * @param {import('vue-i18n').I18n} i18n
 * @returns {import('@vee-validate/i18n').ValidationMessageGenerator}
 */
const createGenerateMessage = (i18n) => {
  const localizedGenerateMessage = localize(veeValidateLocales);

  setFallbackLocale(DEFAULT_LOCALE);
  setLocale(unref(i18n.global.locale) || DEFAULT_LOCALE);

  watch(
    () => unref(i18n.global.locale),
    (locale) => {
      if (locale) {
        setLocale(locale);
      }
    },
  );

  return (ctx) => getCustomRuleMessage(ctx) ?? localizedGenerateMessage(ctx);
};

/**
 * Registers project-specific decimal validation rule.
 *
 * Accepts integer and decimal numbers with `.` or `,` as separator.
 */
const defineDecimalRule = () => {
  defineRule('decimal', (value) => {
    if (isEmpty(value)) {
      return true;
    }

    return /^-?\d+([.,]\d+)?$/.test(String(value));
  });
};

/**
 * Registers project-specific unique validation rule.
 *
 * Rule params:
 * - `validate`: async `(value) => boolean` checker
 * - `values`: array of values that must not contain the current value
 * - `message`: optional custom error message
 */
const defineUniqueRule = () => {
  defineRule('unique', async (value, params) => {
    if (isEmpty(value)) {
      return true;
    }

    if (params?.validate) {
      return await params.validate(value);
    }

    if (params?.values) {
      return !params.values.includes(value);
    }

    return true;
  });
};

/**
 * VeeValidate plugin.
 *
 * Registers global rules from `@vee-validate/rules`, project-specific rules,
 * and configures localized message generation for en/fr locales.
 */
export default {
  /**
   * Install vee-validate plugin.
   *
   * @param {import('vue').App} app
   * @param {Object} [options]
   * @param {import('vue-i18n').I18n} [options.i18n]
   */
  install(app, { i18n } = {}) {
    Object.entries(veeValidateRules).forEach(([name, rule]) => {
      defineRule(name, rule);
    });

    defineDecimalRule();
    defineUniqueRule();

    configure({
      generateMessage: createGenerateMessage(i18n),
    });
  },
};
