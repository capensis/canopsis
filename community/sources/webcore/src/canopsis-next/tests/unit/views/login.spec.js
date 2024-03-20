import { generateRenderer } from '@unit/utils/vue';
import { createAuthModule, createInfoModule, createMockedStoreModules } from '@unit/utils/store';

import Login from '@/views/login.vue';

const stubs = {
  'c-compiled-template': true,
  'login-card': true,
  'login-footer': true,
};

describe('login', () => {
  const { authModule } = createAuthModule();
  const {
    infoModule,
    isCASAuthEnabled,
    isSAMLAuthEnabled,
    isBasicAuthEnabled,
    isOauthAuthEnabled,
    description,
  } = createInfoModule();

  const store = createMockedStoreModules([
    authModule,
    infoModule,
  ]);

  const snapshotFactory = generateRenderer(Login, { stubs });

  it('Renders `login` without auth enabled', () => {
    const wrapper = snapshotFactory({ store });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `login` with all auths enabled', () => {
    isBasicAuthEnabled.mockReturnValueOnce(true);
    isCASAuthEnabled.mockReturnValueOnce(true);
    isSAMLAuthEnabled.mockReturnValueOnce(true);
    isOauthAuthEnabled.mockReturnValueOnce(true);
    description.mockReturnValueOnce('description');

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModule,
        infoModule,
      ]),
    });

    expect(wrapper).toMatchSnapshot();
  });
});
