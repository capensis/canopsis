import { generateRenderer } from '@unit/utils/vue';

import PbehaviorRecurrenceRule from '@/components/other/pbehavior/pbehaviors/partials/pbehavior-recurrence-rule.vue';

const stubs = {
  'recurrence-rule-information': true,
  'pbehavior-recurrence-rule-periods': true,
};

describe('pbehavior-recurrence-rule', () => {
  const pbehavior = {
    _id: 'id-pbehavior',
    rrule: 'RRULE',
  };

  const snapshotFactory = generateRenderer(PbehaviorRecurrenceRule, { stubs });

  test('Renders `pbehavior-recurrence-rule` with pbehavior', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pbehavior,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
