import { generateEntityPatternsTests } from '@unit/utils/patterns';

import DynamicInfoPatternsForm from '@/components/other/dynamic-info/form/fields/dynamic-info-patterns-form.vue';

generateEntityPatternsTests(DynamicInfoPatternsForm, 'dynamic-info-patterns-form', {
  isEntityType: true,
});
