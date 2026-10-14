import { registerHelper, unregisterHelper } from '@/helpers/handlebars';
import { alarmOrEntityTagsHelper } from '@/helpers/handlebars/tags-helper';

/**
 * Mixin that registers the context-aware 'tags' Handlebars helper.
 *
 * This helper renders either c-alarm-tags-chips or c-entity-tags-chips
 * based on the context (alarm or entity).
 */
export const handlebarsTagsHelperMixin = {
  beforeCreate() {
    registerHelper('tags', alarmOrEntityTagsHelper);
  },
  destroyed() {
    unregisterHelper('tags');
  },
};
