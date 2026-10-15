import { generateRenderer } from '@unit/utils/vue';

import ExtraDetailsParents from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-parents.vue';

const stubs = {
  'c-alarm-extra-details-chip': true,
  'c-simple-tooltip': true,
};

describe('extra-details-parents', () => {
  const total = 3;
  const rules = [
    {
      id: 'parent-rule-1-id',
      name: 'parent-rule-1-name',
    },
    {
      id: 'parent-rule-2-id',
      name: 'parent-rule-2-name',
    },
    {
      id: 'parent-rule-3-id',
      name: 'parent-rule-3-name',
    },
  ];

  const snapshotFactory = generateRenderer(ExtraDetailsParents, {
    stubs,
    attachTo: document.body,
  });

  it('Renders `extra-details-parents` with full parents', () => {
    const wrapper = snapshotFactory({
      propsData: {
        total,
        rules,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `extra-details-parents` without rules', () => {
    const wrapper = snapshotFactory({
      propsData: {
        total,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
