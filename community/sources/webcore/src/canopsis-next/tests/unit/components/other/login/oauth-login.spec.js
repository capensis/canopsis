import { generateRenderer } from '@unit/utils/vue';
import { createInfoModule, createMockedStoreModules } from '@unit/utils/store';

import OauthLogin from '@/components/other/login/oauth-login.vue';

describe('oauth-login', () => {
  const { infoModule, oauthConfig } = createInfoModule();
  const store = createMockedStoreModules([infoModule]);

  const snapshotFactory = generateRenderer(OauthLogin);

  it('Renders `oauth-login` with redirect query and one provider', () => {
    oauthConfig.mockReturnValueOnce({ providers: ['google'] });

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([infoModule]),
      mocks: {
        $route: { query: { redirect: '/redirect' } },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `oauth-login` with redirect query and two provider', () => {
    oauthConfig.mockReturnValueOnce({ providers: ['google', 'github'] });

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([infoModule]),
      mocks: {
        $route: { query: { redirect: '/redirect' } },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `oauth-login` without redirect and with two providers', () => {
    oauthConfig.mockReturnValueOnce({ providers: ['google', 'github'] });

    const wrapper = snapshotFactory({
      store,
      mocks: {
        $route: { query: {} },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
