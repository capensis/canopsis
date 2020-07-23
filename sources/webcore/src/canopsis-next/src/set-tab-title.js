import { DEFAULT_APP_TITLE } from '@/config';

/**
 * Function for set tab title
 *
 * @param {String} title
 */
export const setTabTitle = (title = DEFAULT_APP_TITLE) => {
  document.title = title;
};
