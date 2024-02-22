import { generateRenderer } from '@unit/utils/vue';
import { createInfoModule, createMockedStoreModules } from '@unit/utils/store';

import CasLogin from '@/components/other/login/cas-login.vue';

describe('cas-login', () => {
  const { infoModule, casConfig } = createInfoModule();
  const store = createMockedStoreModules([infoModule]);

  const snapshotFactory = generateRenderer(CasLogin);

  it('Renders `cas-login` with redirect query and title', () => {
    casConfig.mockReturnValueOnce({ title: 'Cas config title' });

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([infoModule]),
      mocks: {
        $route: { query: { redirect: '/redirect' } },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `cas-login` without redirect', () => {
    const wrapper = snapshotFactory({
      store,
      mocks: {
        $route: { query: {} },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
