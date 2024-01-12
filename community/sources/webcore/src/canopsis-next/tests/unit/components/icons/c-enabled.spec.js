import { generateRenderer } from '@unit/utils/vue';

import CEnabled from '@/components/icons/c-enabled.vue';

describe('c-enabled', () => {
  const snapshotFactory = generateRenderer(CEnabled);

  it('Renders `c-enabled` correctly.', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-enabled` with enabled prop correctly.', () => {
    const wrapper = snapshotFactory({
      propsData: { value: true },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
