import { omit } from 'lodash';
import { prepareRemediationInstructionsFiltersToQuery } from '@/helpers/filter/remediation-instructions-filter';

export default {
  computed: {
    enabledWidgetRemediationInstructionsFilters() {
      return this.widgetRemediationInstructionsFilters.filter(filter => !filter.disabled);
    },

    remediationInstructionsFilters: {
      get() {
        return this.userPreference.widget_preferences.remediationInstructionsFilters || [];
      },
      set(filters) {
        this.updateWidgetPreferencesInUserPreference({
          ...this.userPreference.widget_preferences,

          remediationInstructionsFilters: filters,
        });

        const newRemediationInstructionsFilters = [
          ...filters,
          ...this.widgetRemediationInstructionsFilters.filter(filter => !filter.disabled),
        ];

        this.updateRemediationInstructionsFiltersInQuery(newRemediationInstructionsFilters);
      },
    },

    widgetRemediationInstructionsFilters: {
      get() {
        const { remediationInstructionsFilters = [] } = this.widget.parameters;
        const { disabledWidgetRemediationInstructionsFilters = [] } = this.userPreference.widget_preferences;

        return remediationInstructionsFilters.map(filter => ({
          ...filter,
          disabled: disabledWidgetRemediationInstructionsFilters.includes(filter._id),
          locked: true,
        }));
      },
      set(filters) {
        this.updateWidgetPreferencesInUserPreference({
          ...this.userPreference.widget_preferences,

          disabledWidgetRemediationInstructionsFilters: filters.filter(filter => filter.disabled)
            .map(filter => filter._id),
        });

        const newRemediationInstructionsFilters = [
          ...this.userPreference.widget_preferences.remediationInstructionsFilters,
          ...filters.filter(filter => !filter.disabled),
        ];

        this.updateRemediationInstructionsFiltersInQuery(newRemediationInstructionsFilters);
      },
    },
  },
  methods: {
    updateRemediationInstructionsFiltersInQuery(filters) {
      const queryWithoutRemediationInstructionsFields = omit(this.query, [
        'with_instructions',
        'include_instructions',
        'exclude_instructions',
        'include_types',
        'exclude_types',
      ]);

      this.query = {
        ...queryWithoutRemediationInstructionsFields,
        ...prepareRemediationInstructionsFiltersToQuery(filters),

        page: 1,
      };
    },
  },
};
