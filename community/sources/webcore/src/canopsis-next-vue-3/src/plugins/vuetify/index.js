import { createVuetify } from 'vuetify';
import { VExpansionPanels } from 'vuetify/components';

import { DEFAULT_THEME_COLORS } from '@/config';

import { createCanopsisVuetifyTheme } from '@/helpers/entities/theme/vuetify';

import { customIconAliases, iconSets } from './icons';

const { definition } = createCanopsisVuetifyTheme(DEFAULT_THEME_COLORS);

export const CANOPSIS_THEME_NAME = 'canopsis';

export const vuetify = createVuetify({
  defaults: {
    VAutocomplete: {
      color: 'primary',
      density: 'compact',
      variant: 'underlined',
    },
    VCombobox: {
      color: 'primary',
      density: 'compact',
      variant: 'underlined',
    },
    VField: {
      color: 'primary',
    },
    VMenu: {
      transition: 'fade-transition',
    },
    VSelect: {
      color: 'primary',
      density: 'compact',
      variant: 'underlined',
      transition: 'fade-transition',
    },
    VTextarea: {
      color: 'primary',
      density: 'compact',
      variant: 'underlined',
    },
    VTextField: {
      color: 'primary',
      density: 'compact',
      variant: 'underlined',
    },
    VTooltip: {
      transition: 'fade-transition',
    },
    VBtn: {
      density: 'default',
    },
    VCard: {
      elevation: 2,
    },
    VExpansionPanels: {
      elevation: 2,
    },
    VTabs: {
      sliderColor: 'primary',
    },
  },
  icons: {
    defaultSet: 'md',
    aliases: customIconAliases,
    sets: iconSets,
  },
  theme: {
    defaultTheme: CANOPSIS_THEME_NAME,
    themes: {
      [CANOPSIS_THEME_NAME]: definition,
    },
  },
});

/**
 * Kebab-case aliases for Vuetify components used as dynamic tags in vuedraggable.
 * vuedraggable resolves `tag` via resolveComponent(string), which requires global registration.
 */
export const registerVuetifyComponentAliases = (app) => {
  app.component('v-expansion-panels', VExpansionPanels);
};

export default vuetify;
