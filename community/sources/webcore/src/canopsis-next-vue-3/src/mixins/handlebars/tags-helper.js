import { registerHelper, unregisterHelper } from '@/helpers/handlebars';
import { alarmOrEntityTagsHelper } from '@/helpers/handlebars/tags-helper';

/**
 * Registers the context-aware tags Handlebars helper while the component is mounted.
 */
export const handlebarsTagsHelperMixin = {
  beforeCreate() {
    registerHelper('tags', alarmOrEntityTagsHelper);
  },

  beforeUnmount() {
    unregisterHelper('tags');
  },
};
