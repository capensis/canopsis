import { generateRenderer } from '@unit/utils/vue';
import { createAuthModule, createInfoModule, createMockedStoreModules } from '@unit/utils/store';

import Login from '@/views/login.vue';

const stubs = {
  'c-compiled-template': true,
  'base-login': true,
  'cas-login': true,
  'saml-login': true,
  'login-footer': true,
};

describe('login', () => {
  const { authModule } = createAuthModule();
  const { infoModule, isCASAuthEnabled, isSAMLAuthEnabled, description } = createInfoModule();
  const store = createMockedStoreModules([
    authModule,
    infoModule,
  ]);

  const snapshotFactory = generateRenderer(Login, { stubs });

  it('Renders `login` with default state', () => {
    const wrapper = snapshotFactory({ store });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `login` with custom state', () => {
    isCASAuthEnabled.mockReturnValueOnce(true);
    isSAMLAuthEnabled.mockReturnValueOnce(true);
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
