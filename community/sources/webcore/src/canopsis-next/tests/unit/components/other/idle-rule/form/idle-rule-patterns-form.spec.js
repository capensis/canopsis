import { generateEntityPatternsTests } from '@unit/utils/patterns';

import IdleRulePatternsForm from '@/components/other/idle-rule/form/idle-rule-patterns-form.vue';

generateEntityPatternsTests(IdleRulePatternsForm, 'idle-rule-patterns-form', { isEntityType: true });
