import { AVAILABILITY_GROUPED_ADVANCED_SEARCH_FIELDS } from '@/constants';

import { useEntityAdvancedSearchAttributes } from './attributes-map-entity';

export const useAvailabilityAdvancedSearchAttributes = () => useEntityAdvancedSearchAttributes({
  groups: AVAILABILITY_GROUPED_ADVANCED_SEARCH_FIELDS,
});
