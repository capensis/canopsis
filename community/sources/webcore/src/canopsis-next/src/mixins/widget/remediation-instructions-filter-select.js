import { get, omit } from 'lodash';
import { prepareRemediationInstructionsFiltersToQuery } from '@/helpers/filter/remediation-instructions-filter';

export default {
  computed: {
    remediationInstructionsFilters: {
      get() {
        return this.userPreference.content.remediationInstructionsFilters || [];
      },
      set(filters) {
        this.updateContentInUserPreference({
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
        const { disabledWidgetRemediationInstructionsFilters = [] } = this.userPreference.content;

        return remediationInstructionsFilters.map(filter => ({
          ...filter,
          disabled: disabledWidgetRemediationInstructionsFilters.includes(filter._id),
          locked: true,
        }));
      },
      set(filters) {
        this.updateContentInUserPreference({
          disabledWidgetRemediationInstructionsFilters: filters.filter(filter => filter.disabled)
            .map(filter => filter._id),
        });

        const newRemediationInstructionsFilters = [
          ...get(this.userPreference, 'content.remediationInstructionsFilters', []),
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
        'include_instruction_types',
        'exclude_instruction_types',
      ]);

      this.query = {
        ...queryWithoutRemediationInstructionsFields,
        ...prepareRemediationInstructionsFiltersToQuery(filters),

        page: 1,
      };
    },
  },
};
