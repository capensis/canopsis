import { omit } from 'lodash';

import { normalizeHtml } from '@/helpers/html';

export const vuetifyCustomIconsBaseMixin = {
  methods: {
    hasIconInVuetify(name) {
      return !!this.$vuetify.icons.values[name];
    },

    registerIconInVuetify(name, template, meta = {}) {
      this.$vuetify.icons.values = {
        ...this.$vuetify.icons.values,

        [name]: {
          meta,
          component: {
            template: normalizeHtml(template),
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
