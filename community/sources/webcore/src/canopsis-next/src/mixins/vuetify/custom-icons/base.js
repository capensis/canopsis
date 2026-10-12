import {
  customIconsRegistry,
  registerCustomIcon,
  unregisterCustomIcon,
} from '@/plugins/vuetify/icons';

import { normalizeHtml } from '@/helpers/html';

export const vuetifyCustomIconsBaseMixin = {
  methods: {
    hasIconInVuetify(name) {
      return !!customIconsRegistry[name];
    },

    registerIconInVuetify(name, template, meta = {}) {
      registerCustomIcon(name, normalizeHtml(template), meta);
    },

    unregisterIconFromVuetify(name) {
      if (!this.hasIconInVuetify(name)) {
        return;
      }

      unregisterCustomIcon(name);
    },
  },
};
