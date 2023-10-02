import { generateEntityPatternsTests } from '@unit/utils/patterns';

import ScenarioActionPatternsForm from '@/components/other/scenario/form/scenario-action-patterns-form.vue';

generateEntityPatternsTests(ScenarioActionPatternsForm, 'scenario-action-patterns-form', {
  name: 'customName',
});
