import { generateRenderer } from '@unit/utils/vue';

import ExtraDetailsChildren from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-children.vue';

const stubs = {
  'c-alarm-extra-details-chip': true,
  'c-simple-tooltip': true,
};

describe('extra-details-children', () => {
  const total = 3;
  const opened = 2;
  const closed = 1;
  const rule = {
    name: 'rule-name',
  };

  const snapshotFactory = generateRenderer(ExtraDetailsChildren, {
    stubs,
    attachTo: document.body,
  });

  it('Renders `extra-details-children` with full children and rule', () => {
    jest.useFakeTimers();

    const wrapper = snapshotFactory({
      propsData: {
        total,
        opened,
        closed,
        rule,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
