import { inject, set, onBeforeUnmount } from 'vue';

import { PATTERNS_FIELDS } from '@/constants';

const PATTERN_FIELDS_TO_COMPONENTS_CLASSES = {
  [PATTERNS_FIELDS.alarm]: 'c-alarm-patterns-field',
  [PATTERNS_FIELDS.entity]: 'c-entity-patterns-field',
  [PATTERNS_FIELDS.pbehavior]: 'c-pbehavior-patterns-field',
  [PATTERNS_FIELDS.event]: 'c-event-filter-patterns-field',
  [PATTERNS_FIELDS.totalEntity]: 'c-entity-patterns-field',
  [PATTERNS_FIELDS.serviceWeather]: 'c-service-weather-patterns-field',
};

export const PATTERNS_FIELDS_TO_EXPANDED_KEYS = {
  [PATTERNS_FIELDS.alarm]: 'alarm',
  [PATTERNS_FIELDS.entity]: 'entity',
  [PATTERNS_FIELDS.pbehavior]: 'pbehavior',
  [PATTERNS_FIELDS.event]: 'event',
  [PATTERNS_FIELDS.totalEntity]: 'totalEntity',
  [PATTERNS_FIELDS.serviceWeather]: 'serviceWeather',
};

/**
 * Resolves when the Vuetify expansion panel content `transitionend` on `parentElement` has fired,
 * so scrolling runs after the expand/collapse animation.
 *
 * @param {HTMLElement} parentElement - Container that includes `.v-expansion-panel-content`
 */
const waitCollapsePanelTransitionEnd = (parentElement) => {
  const element = parentElement.querySelector('.v-expansion-panel-content');

  if (!element) {
    return Promise.resolve();
  }

  return new Promise((resolve) => {
    const onEnd = (event) => {
      if (event.target !== element) return;

      element.removeEventListener('transitionend', onEnd);
      resolve();
    };

    element.addEventListener('transitionend', onEnd);
  });
};

/**
 * Waits for any expansion transition on the element, then scrolls it into view with smooth behavior.
 *
 * @param {HTMLElement} element - Target node (e.g. a `.c-collapse-panel` root)
 */
const scrollToElement = async (element) => {
  await waitCollapsePanelTransitionEnd(element);

  element.scrollIntoView({ behavior: 'smooth' });
};

/**
 * Registers expand-and-scroll behavior with AI chat so the correct pattern sub-panels open and the view follows.
 *
 * @param {Object} options
 * @param {import('vue').Ref} options.wrapperElement - Ref to the field wrapper that contains pattern sections
 * @param {import('vue').Ref} options.expanded - Ref to expanded flags for each pattern subsection
 */
export const usePatternsFieldExpand = ({ wrapperElement, expanded }) => {
  const aiChat = inject('$aiChat', {});

  /**
   * Expands the pattern sub-panels for the given field keys, then scrolls to the first field's section.
   * Used e.g. by AI chat to focus the user on relevant pattern blocks.
   *
   * @param {Object} args
   * @param {string[]} args.fields - Pattern field names (e.g. alarm, entity)
   */
  const expandPatternsByFields = async ({ fields }) => {
    const firstClass = PATTERN_FIELDS_TO_COMPONENTS_CLASSES[fields[0]];

    fields.forEach((field) => {
      const expandedKey = PATTERNS_FIELDS_TO_EXPANDED_KEYS[field];

      if (!expanded.value[expandedKey]) {
        set(expanded.value, expandedKey, true);
      }
    });

    scrollToElement(wrapperElement.value.querySelector(`.${firstClass}`).closest('.c-collapse-panel'));
  };

  aiChat?.registerExpandFunction?.(expandPatternsByFields);

  onBeforeUnmount(() => aiChat?.unregisterExpandFunction?.(expandPatternsByFields));
};
