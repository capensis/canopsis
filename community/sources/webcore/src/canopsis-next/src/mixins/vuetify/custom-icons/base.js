import {
  customIconsRegistry,
  registerCustomIcon,
  unregisterCustomIcon,
} from '@/plugins/vuetify/icons';

export const vuetifyCustomIconsBaseMixin = {
  methods: {
    hasIconInVuetify(name) {
      return !!customIconsRegistry[name];
    },

    registerIconInVuetify(name, template, meta = {}) {
      registerCustomIcon(name, template, meta);
    },

    unregisterIconFromVuetify(name) {
      if (!this.hasIconInVuetify(name)) {
        return;
      }

      unregisterCustomIcon(name);
    },
  },
};
