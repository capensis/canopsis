import Faker from 'faker';

import { flushPromises, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createAuthModule, createInfoModule, createMockedStoreModules } from '@unit/utils/store';
import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { mockRouter } from '@unit/utils/mock-hooks';

import { ROUTES_NAMES } from '@/constants';

import BasicLogin from '@/components/other/login/basic-login.vue';

const stubs = {
  'ldap-login-information': true,
  'login-form': true,
  'c-compiled-template': true,
  'c-alert': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const snapshotStubs = {
  'ldap-login-information': true,
  'login-form': true,
  'c-alert': true,
  'c-compiled-template': true,
};

const selectSubmitButton = wrapper => wrapper.find('button.v-btn');
const selectLoginForm = wrapper => wrapper.find('login-form-stub');
const selectAlert = wrapper => wrapper.find('c-alert-stub');

describe('basic-login', () => {
  const $router = mockRouter();

  const { authModule, currentUser, login } = createAuthModule();
  const { infoModule, isLDAPAuthEnabled, footer } = createInfoModule();
  const store = createMockedStoreModules([
    authModule,
    infoModule,
  ]);

  const factory = generateShallowRenderer(BasicLogin, {
    stubs,
    attachTo: document.body,
    mocks: { $router },
  });
  const snapshotFactory = generateRenderer(BasicLogin, { stubs: snapshotStubs });

  it('Form submitted after trigger submit button', async () => {
    const wrapper = factory({
      store,
      mocks: {
        $route: { query: {} },
        $router: {},
      },
    });

    const username = Faker.datatype.string();
    const password = Faker.datatype.string();

    selectLoginForm(wrapper).triggerCustomEvent('input', username, ['username']);
    selectLoginForm(wrapper).triggerCustomEvent('input', password, ['password']);

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(login).toBeCalledWith(
      expect.any(Object),
      { username, password },
      undefined,
    );

    expect($router.push).toBeCalledWith({ name: ROUTES_NAMES.home });
  });

  it('Error shown after trigger submit button with error', async () => {
    login.mockRejectedValueOnce({});

    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        infoModule,
      ]),
      mocks: {
        $route: { query: {} },
        $router: {},
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(login).toBeCalledWith(
      expect.any(Object),
      { username: '', password: '' },
      undefined,
    );

    expect(selectAlert(wrapper).attributes('value')).toBeTruthy();

    expect($router.push).not.toBeCalled();
  });

  it('Form submitted after trigger submit button with redirect query', async () => {
    const redirectUrl = '/redirect/url';
    const wrapper = factory({
      store,
      mocks: {
        $route: { query: { redirect: redirectUrl } },
        $router: {},
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(login).toBeCalledWith(
      expect.any(Object),
      {
        username: '',
        password: '',
      },
      undefined,
    );

    expect($router.push).toBeCalledWith(redirectUrl);
  });

  it('Form submitted after trigger submit button with default view', async () => {
    const defaultView = {
      _id: Faker.datatype.string(),
    };
    currentUser.mockReturnValueOnce({
      defaultview: defaultView,
    });
    const wrapper = factory({
      store: createMockedStoreModules([
        authModule,
        infoModule,
      ]),
      mocks: {
        $route: { query: {} },
        $router: {},
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(login).toBeCalledWith(
      expect.any(Object),
      {
        username: '',
        password: '',
      },
      undefined,
    );

    expect($router.push).toBeCalledWith({
      name: ROUTES_NAMES.view,
      params: { id: defaultView._id },
    });
  });

  it('Renders `base-login` with default state', () => {
    const wrapper = snapshotFactory({ store });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `base-login` with custom state', () => {
    isLDAPAuthEnabled.mockReturnValueOnce(true);
    footer.mockReturnValueOnce('footer');

    const wrapper = snapshotFactory({
      store: createMockedStoreModules([
        authModule,
        infoModule,
      ]),
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `base-login` with error', async () => {
    const wrapper = snapshotFactory({
      store,
      mocks: {
        $route: { query: {} },
        $router: {},
      },
    });

    await wrapper.setData({ serverErrorMessage: 'Server error message' });

    expect(wrapper).toMatchSnapshot();
  });
});
