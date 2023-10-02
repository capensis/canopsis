import { generateEntityPatternsTests } from '@unit/utils/patterns';

import PbehaviorPatternsForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-patterns-form.vue';

generateEntityPatternsTests(PbehaviorPatternsForm, 'pbehavior-patterns-form', { readonly: true });
