import { generateRenderer } from '@unit/utils/vue';

import { DENSE_TYPES } from '@/constants';

import CDensityBtnToggle from '@/components/common/groups/c-density-btn-toggle.vue';

describe('c-density-btn-toggle', () => {
  const snapshotFactory = generateRenderer(CDensityBtnToggle);

  it('Renders `c-density-btn-toggle` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-density-btn-toggle` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: DENSE_TYPES.medium,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
