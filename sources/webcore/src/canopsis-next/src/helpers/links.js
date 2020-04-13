import { isObject } from 'lodash';

import { getCollectionComparator } from './sort';

/**
 * The linkbuilders used to return the links directly as
 * strings. They can now also return objects with the
 * properties 'label' and 'link', allowing to change the link's
 * label.
 * The following code converts the "legacy" representation
 * (strings) into the "new" representation, so they can be
 * displayed in the same manner by the template.
 *
 * @param {Array} links
 * @param {string} category
 * @return {{ label: string, link: string }[]}
 */
export function harmonizeLinks(links, category) {
  return links.map((link, index) => {
    if (isObject(link) && link.link && link.label) {
      return link;
    }

    return {
      label: `${category} - ${index}`,
      link,
    };
  }).sort(getCollectionComparator('label'));
}

/**
 * Category harmonization
 *
 * @param {Object} categories
 * @param {Function} additionalFilter
 * @return {Array}
 */
export function harmonizeCategories(categories, additionalFilter = () => true) {
  return Object.entries(categories).reduce((acc, [category, categoryLinks]) => {
    if (categoryLinks.length && additionalFilter(category)) {
      acc.push({
        label: category,
        links: harmonizeLinks(categoryLinks, category),
      });
    }

    return acc;
  }, []).sort(getCollectionComparator('label'));
}
