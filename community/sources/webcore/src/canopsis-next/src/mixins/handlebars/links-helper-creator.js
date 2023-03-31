import { get } from 'lodash';
import Handlebars from 'handlebars';

import { registerHelper, unregisterHelper } from '@/helpers/handlebars';

import { authMixin } from '../auth';

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
      const category = hash.category ? `'${hash.category}'` : undefined;

      return new Handlebars.SafeString(`
          <c-links-list v-if="links && hasAccessForLinks" :links="links" :category="${category}" />
        `);
    });
  },
  destroyed() {
    unregisterHelper('links');
  },
});
