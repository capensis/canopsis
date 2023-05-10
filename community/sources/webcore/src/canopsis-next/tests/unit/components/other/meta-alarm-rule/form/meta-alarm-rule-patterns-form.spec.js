import { generateEntityPatternsTests } from '@unit/utils/patterns';

import MetaAlarmRulePatternsForm from '@/components/other/meta-alarm-rule/form/meta-alarm-rule-patterns-form.vue';

generateEntityPatternsTests(MetaAlarmRulePatternsForm, 'meta-alarm-rule-patterns-form', { disabled: true });
