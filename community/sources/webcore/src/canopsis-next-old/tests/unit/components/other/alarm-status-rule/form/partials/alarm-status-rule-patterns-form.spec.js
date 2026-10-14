import { generateEntityPatternsTests } from '@unit/utils/patterns';

import AlarmStatusRulePatternsForm from '@/components/other/alarm-status-rule/form/alarm-status-rule-patterns-form.vue';

generateEntityPatternsTests(AlarmStatusRulePatternsForm, 'alarm-status-rule-patterns-form', {
  flapping: true,
  disabled: true,
});
