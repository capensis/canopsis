import { generateRenderer } from '@unit/utils/vue';

import ImpactStateIndicator from '@/components/widgets/service-weather/impact-state-indicator.vue';

describe('impact-state-indicator', () => {
  const snapshotFactory = generateRenderer(ImpactStateIndicator);

  it('Renders `impact-state-indicator` with default props correctly', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `impact-state-indicator` with custom props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 8,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
