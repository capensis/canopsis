/**
 * Remove trailing slashes from url (http://example.com//login -> http://example.com/login)
 *
 * @param {string} [url = '']
 * @returns {string}
 */
export const removeTrailingSlashes = (url = '') => url.replace(/([^:]\/)\/+/g, '$1');

/**
 * We need to use this function to avoid problem with double slashes in url
 *
 * @example http://example.com//something//something will be redirected to http://example.com/something/something
 */
export const reloadPageWithTrailingSlashes = () => {
  const {
    origin,
    pathname,
    search,
    href,
  } = window.location;
  const preparedHref = `${removeTrailingSlashes(`${origin}${pathname}`)}${search}`;

  if (href !== preparedHref) {
    window.location = preparedHref;
  }
};

/**
 * Open url in new tab
 *
 * @param url
 * @return {WindowProxy}
 */
export const openUrlInNewTab = url => window.open(removeTrailingSlashes(url), '_blank');
