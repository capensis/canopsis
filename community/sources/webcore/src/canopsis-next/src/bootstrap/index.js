import deepFreeze from 'deep-freeze';
import WebFont from 'webfontloader';

import * as constants from '@/constants';
import * as config from '@/config';

/* eslint-disable import/order */
import { createVuetify } from '@/plugins/vuetify';
import { bootstrapApplicationPlugins } from './plugins';
import { registerApplicationComponents } from './components';
/* eslint-enable import/order */

import router from '@/router';
import store from '@/store';
import i18n from '@/i18n';

import { themePropertiesToCSSVariables } from '@/helpers/entities/theme/entity';

import App from '@/app.vue';

/**
 * @param {import('vue').VueConstructor | import('vue').Vue} Vue
 * @returns {*}
 */
export const bootstrapApplication = (Vue) => {
  WebFont.load({
    custom: {
      families: ['Roboto:300,400,500,700', 'Material+Icons'],
      urls: ['/styles/fonts.css'],
    },
  });

  const vuetify = createVuetify(Vue, {
    theme: {
      dark: false,
      themes: {
        light: themePropertiesToCSSVariables(config.DEFAULT_THEME_COLORS),
        dark: themePropertiesToCSSVariables(config.DEFAULT_THEME_COLORS),
      },
    },
  });
  bootstrapApplicationPlugins(Vue);
  registerApplicationComponents(Vue);

  Vue.config.productionTip = false;

  Vue.config.errorHandler = (err) => {
    console.error(err);

    store.dispatch('popups/error', { text: err.description || i18n.t('errors.default') });
  };

  if (process.env.NODE_ENV === 'development') {
    Vue.config.devtools = true;
    Vue.config.performance = true;
  }

  Vue.prototype.$constants = deepFreeze(constants);
  Vue.prototype.$config = deepFreeze(config);

  return new Vue({
    vuetify,
    router,
    store,
    i18n,
    render: h => h(App),
  });
};
