import { escapeRegExp } from 'lodash';
import Handlebars from 'handlebars';

import { MAX_LIMIT } from '@/constants';

/**
 * Context-aware Handlebars helper that renders tags chips component based on the available context.
 *
 * Returns c-alarm-tags-chips if alarm is present in the context.
 * Returns c-entity-tags-chips if only entity is present in the context.
 * Returns empty string if neither alarm nor entity is available.
 *
 * @returns {Handlebars.SafeString} The rendered tags chips component
 *
 * @example
 * // In a template with alarm context:
 * {{tags}} // Renders c-alarm-tags-chips
 *
 * @example
 * // In a template with entity context:
 * {{tags}} // Renders c-entity-tags-chips
 */
export const alarmOrEntityTagsHelper = function tagsHelper(...args) {
  const { alarm, entity } = this;

  if (alarm) {
    const { hash: { regex = '', flags = '' } = {} } = args.pop() ?? {};
    const escapedArgs = args.filter(Boolean).map(arg => escapeRegExp(arg));
    const nameFilter = escapedArgs.length ? `^(${escapedArgs.join('|')})$` : '';

    return new Handlebars.SafeString(
      `<c-alarm-tags-chips
        :alarm="alarm"
        :selected-tags="selectedTags"
        name-filter="${nameFilter}"
        regex-filter="${regex}"
        regex-filter-flags="${flags}"
        inline-count="${MAX_LIMIT}"
        closable-active
        @select="$emit('select:tag', $event)"
        @close="$emit('remove:tag', $event)"
      ></c-alarm-tags-chips>`,
    );
  }

  if (entity) {
    return new Handlebars.SafeString(
      `<c-entity-tags-chips
        :entity="entity"
        inline-count="${MAX_LIMIT}"
       ></c-entity-tags-chips>`,
    );
  }

  return new Handlebars.SafeString('');
};
