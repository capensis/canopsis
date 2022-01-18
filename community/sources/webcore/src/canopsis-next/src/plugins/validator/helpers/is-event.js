import { isCallable } from './is-callable';

/**
 * Check is validate event
 *
 * @param {any} evt
 * @return {boolean}
 */
export const isEvent = evt => (
  typeof Event !== 'undefined'
  && isCallable(Event)
  && evt instanceof Event
) || (evt && evt.srcElement);
