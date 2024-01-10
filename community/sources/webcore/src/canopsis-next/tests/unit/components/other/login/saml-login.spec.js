import { generateRenderer } from '@unit/utils/vue';
import { createInfoModule, createMockedStoreModules } from '@unit/utils/store';

import SamlLogin from '@/components/other/login/saml-login.vue';

describe('saml-login', () => {
  const { infoModule, samlConfig } = createInfoModule();
  const store = createMockedStoreModules([infoModule]);

  const snapshotFactory = generateRenderer(SamlLogin);

  it('Renders `saml-login` with redirect query and title', () => {
    samlConfig.mockReturnValueOnce({ title: 'Saml config title' });

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([infoModule]),
      mocks: {
        $route: { query: { redirect: '/redirect' } },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `saml-login` without redirect', () => {
    const wrapper = snapshotFactory({
      store,
      mocks: {
        $route: { query: {} },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
