import { get } from 'lodash';
import Handlebars from 'handlebars';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars';

import { authMixin } from '../auth';

/**
 * Create a mixin that registers the links helper with fixed Vue bindings.
 * Handlebars resolves the category before it is escaped as a static attribute.
 *
 * @param {string} linksPath - Path to the links data on the component instance.
 * @param {string} right - Permission required to display links.
 * @returns {Object} Vue mixin managing helper registration and link access.
 */
export const handlebarsLinksHelperCreator = (linksPath, right) => ({
  mixins: [authMixin],
  computed: {
    hasAccessForLinks() {
      return this.checkAccess(right);
    },

    links() {
      return get(this, linksPath);
    },
  },
  beforeCreate() {
    registerHelper('links', ({ hash }) => {
      const category = hash.category
        ? `category="${Handlebars.escapeExpression(String(hash.category))}"`
        : '';

      return new Handlebars.SafeString(`
          <c-links-list v-if="links && hasAccessForLinks" :links="links" ${category} />
        `);
    });
  },
  destroyed() {
    unregisterHelper('links');
  },
});
