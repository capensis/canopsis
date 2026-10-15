import { generateRenderer } from '@unit/utils/vue';
import { createInfoModule, createMockedStoreModules } from '@unit/utils/store';

import LoginFooter from '@/components/other/login/login-footer.vue';

describe('login-footer', () => {
  const { infoModule, version } = createInfoModule();
  const store = createMockedStoreModules([infoModule]);

  const snapshotFactory = generateRenderer(LoginFooter);

  it('Renders `login-footer` without version', () => {
    const wrapper = snapshotFactory({ store });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `login-footer` without version', () => {
    version.mockReturnValueOnce('23.10');

    const wrapper = snapshotFactory({ store: createMockedStoreModules([infoModule]) });

    expect(wrapper).toMatchSnapshot();
  });
});
