import { generateRenderer } from '@unit/utils/vue';

import CProgressOverlay from '@/components/common/overlay/c-progress-overlay.vue';

describe('c-progress-overlay', () => {
  const snapshotFactory = generateRenderer(CProgressOverlay);

  it('Renders `c-progress-overlay` with default props correctly', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-progress-overlay` with pending true correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pending: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-progress-overlay` with custom props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        pending: true,
        opacity: 0.3,
        backgroundColor: 'black',
        color: 'secondary',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
