import { isArray } from 'lodash';

import { setField } from '@/helpers/immutable';

export const sideBarSettingsMainFilterIntersection = {
  methods: {
    /**
     * Get prepared userPreferences for request sending
     *
     * @returns {Object}
     */
    getPreparedUserPreference() {
      const getActualFilter = (filter) => {
        if (!filter) {
          return filter;
        }

        const intersectFilter = this.settings.widget.parameters.viewFilters
          .find(viewFilter => viewFilter.title === filter.title);

        if (intersectFilter) {
          return intersectFilter;
        }

        return filter;
      };

      const { mainFilter } = this.userPreference.content;
      const newMainFilter = isArray(mainFilter)
        ? mainFilter.map(getActualFilter)
        : getActualFilter(mainFilter);

      return setField(this.userPreference, 'content', value => ({
        ...value,
        ...this.settings.userPreferenceContent,

        mainFilter: newMainFilter,
      }));
    },
  },
};
