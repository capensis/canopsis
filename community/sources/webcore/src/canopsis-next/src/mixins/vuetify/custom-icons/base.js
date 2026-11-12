import { omit } from 'lodash';

import { sanitizeSvg } from '@/helpers/html';

export const vuetifyCustomIconsBaseMixin = {
  methods: {
    hasIconInVuetify(name) {
      return !!this.$vuetify.icons.values[name];
    },

    /**
     * Register an icon using sanitized SVG data inside a fixed Vue wrapper.
     *
     * @param {string} name - Icon name used by the Vuetify registry.
     * @param {string} template - Uploaded or previously stored SVG markup.
     * @param {Object} [meta={}] - Icon metadata preserved in the registry.
     * @param {boolean} [meta.fillBorder] - Whether the icon fills its chip border.
     * @returns {void}
     */
    registerIconInVuetify(name, template, meta = {}) {
      const sanitized = sanitizeSvg(template);
      const element = document.createElement('template');
      element.innerHTML = sanitized;
      const svg = element.content.querySelector('svg');
      const attributes = svg ? Object.fromEntries([...svg.attributes].map(({ name: key, value }) => [key, value])) : {};
      this.$vuetify.icons.values = {
        ...this.$vuetify.icons.values,

        [name]: {
          meta,
          component: {
            /**
             * Compile only this fixed wrapper. Uploaded SVG attributes are data,
             * and its children enter the DOM through v-html, never Vue.compile.
             */
            template: '<svg v-bind="svgAttributes" v-html="svgContent"></svg>',
            data: () => ({
              svgAttributes: attributes,
              svgContent: svg?.innerHTML || '',
            }),
          },
        },
      };
    },

    unregisterIconFromVuetify(name) {
      if (!this.hasIconInVuetify(name)) {
        return;
      }

      this.$vuetify.icons.values = omit(this.$vuetify.icons.values, [name]);
    },
  },
};
