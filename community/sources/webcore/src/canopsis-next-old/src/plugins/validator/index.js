import { isObject } from 'lodash';
import VeeValidate, { Validator, Rules, ErrorBag } from 'vee-validate';

import { isValidJson } from './helpers/is-valid-json';
import { isValidUrl } from './helpers/is-valid-url';
import { isUniqueValue } from './helpers/is-unique-value';
import { debounce } from './helpers/debounce';
import { isEvent } from './helpers/is-event';
import { findField } from './helpers/find-field';
import { isValidPicker } from './helpers/is-valid-picker';

const REQUIRED_RULE = 'required';

const getParentValidatorOptions = (vnode) => {
  const validateOptions = vnode.$options.$_veeValidate;

  if (validateOptions?.validator === 'new') {
    return validateOptions;
  }

  if (!vnode?.$parent) {
    return {};
  }

  return getParentValidatorOptions(vnode.$parent);
};

/**
 * Required fields are marked by asterisk and error color only — hide the required rule message.
 *
 * @param {Object} error
 * @returns {string}
 */
const getDisplayErrorMessage = error => (error.rule === REQUIRED_RULE ? '' : error.msg);

/**
 * Maps ErrorBag collect result to display messages, blanking required rule texts.
 *
 * @param {Array|Object} collected
 * @returns {Array|Object}
 */
const collectedToDisplayMessages = (collected) => {
  if (!Array.isArray(collected)) {
    return Object.keys(collected).reduce((acc, key) => {
      acc[key] = collectedToDisplayMessages(collected[key]);

      return acc;
    }, {});
  }

  return collected.map(getDisplayErrorMessage);
};

Validator.prototype.remove = (name) => {
  delete Rules[name];
};

const sourceErrorBagCollect = ErrorBag.prototype.collect;

/**
 * Keep required errors for state (`has`, border color), but do not expose their messages in UI.
 */
ErrorBag.prototype.collect = function collect(field, scope, map = true) {
  const collected = sourceErrorBagCollect.call(this, field, scope, false);

  if (!map) {
    return collected;
  }

  return collectedToDisplayMessages(collected);
};

ErrorBag.prototype.has = function has(field, scope) {
  return this.collect(field, scope, false).length > 0;
};

ErrorBag.prototype.first = function first(field, scope) {
  return this.collect(field, scope)[0];
};

export default {
  install(Vue, { i18n } = {}) {
    Vue.use(VeeValidate, {
      i18n,
      inject: false,
      aria: true,
      silentTranslationWarn: false,
    });

    Validator.extend('json', {
      getMessage: () => i18n.t('errors.JSONNotValid'),
      validate: isValidJson,
    });

    Validator.extend('unique', {
      getMessage: () => i18n.t('errors.unique'),
      validate: isUniqueValue,
    }, {
      paramNames: ['values', 'initialValue'],
    });

    Validator.extend('picker_format', {
      getMessage: () => i18n.t('errors.endDateLessOrEqualStartDate'),
      validate: isValidPicker,
    }, {
      paramNames: ['preparer'],
    });

    Validator.extend('picker_format', {
      getMessage: () => i18n.t('errors.endDateLessOrEqualStartDate'),
      validate: isValidPicker,
    }, {
      paramNames: ['preparer'],
    });

    Validator.extend('url', { validate: isValidUrl });

    const sourceDirective = Vue.directive('validate');

    Vue.directive('validate', {
      ...sourceDirective,

      /* eslint-disable */
      bind(el, binding, vnode) {
        const validatorOptions = getParentValidatorOptions(vnode.context);

        sourceDirective.bind.call(this, el, binding, vnode);

        const field = findField(el, vnode.context);

        if (field && validatorOptions.delay) {
          field.delay = validatorOptions.delay;
        }

        if (field && isObject(field.initialValue)) {
          field.__proto__.addValueListeners = function addValueListeners() {
            this.unwatch(/^input_.+/);
            if (!this.listen || !this.el) return;

            const token = { cancelled: false };
            const fn = this.targetOf
              ? () => {
                const target = this.validator._resolveField(`#${this.targetOf}`);
                if (target && target.flags.validated) {
                  this.validator.validate(`#${this.targetOf}`);
                }
              }
              : (...args) => {
              // if its a DOM event, resolve the value, otherwise use the first parameter as the value.
                if (args.length === 0 || isEvent(args[0])) {
                  args[0] = this.value;
                }

                /* 'We've replaced `addValueListeners` only for non primitives
                 * because we don't need run validation when object was not changed
                 * (For example: in case, when parent object was change by updating another field.)'
                 */
                if (args[0] === args[1]) {
                  return;
                }

                this.flags.pending = true;
                this._cancellationToken = token;
                this.validator.validate(`#${this.id}`, args[0]);
              };

            const inputEvent = this._determineInputEvent();
            let events = this._determineEventList(inputEvent);

            // if on input validation is requested.
            if (events.includes(inputEvent)) {
              let ctx = null;
              let expression = null;
              let watchCtxVm = false;
              // if its watchable from the context vm.
              if (this.model && this.model.expression) {
                ctx = this.vm;
                expression = this.model.expression;
                watchCtxVm = true;
              }

              // watch it from the custom component vm instead.
              if (!expression && this.componentInstance && this.componentInstance.$options.model) {
                ctx = this.componentInstance;
                expression = this.componentInstance.$options.model.prop || 'value';
              }

              if (ctx && expression) {
                const debouncedFn = debounce(fn, this.delay[inputEvent], token);
                const unwatch = ctx.$watch(expression, debouncedFn);
                this.watchers.push({
                  tag: 'input_model',
                  unwatch: () => {
                    this.vm.$nextTick(() => {
                      unwatch();
                    });
                  },
                });

                // filter out input event when we are watching from the context vm.
                if (watchCtxVm) {
                  events = events.filter(e => e !== inputEvent);
                }
              }
            }

            // Add events.
            events.forEach((e) => {
              const debouncedFn = debounce(fn, this.delay[e], token);

              this._addComponentEventListener(e, debouncedFn);
              this._addHTMLEventListener(e, debouncedFn);
            });
          };

          field.addValueListeners();
        }
      },
      /* eslint-enable */
    });
  },
};
