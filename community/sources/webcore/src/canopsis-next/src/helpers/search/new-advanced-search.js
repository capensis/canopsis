import { patternRuleToForm } from '@/helpers/entities/pattern/form';

export const advancedSearchRuleToForm = (advancedSearch = {}) => ({
  ...patternRuleToForm(advancedSearch),
  rangeValue: {
    from: null,
    to: null,
  },
  filled: [],
});
