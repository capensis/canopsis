import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { mockRouter } from '@unit/utils/mock-hooks';
import { createAuthModule, createInfoModule, createMockedStoreModules } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

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

const selectSubmitButton = wrapper => wrapper.find('button.v-btn');
const selectLoginForm = wrapper => wrapper.find('login-form-stub');
const selectAlert = wrapper => wrapper.find('c-alert-stub');

describe('basic-login', () => {
  const $router = mockRouter();

  const { authModule, currentUser, login } = createAuthModule();
  const { infoModule } = createInfoModule();
  const store = createMockedStoreModules([authModule, infoModule]);

  const factory = generateShallowRenderer(BasicLogin, {
    stubs,
    attachTo: document.body,
    mocks: { $router },
  });

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

    expect(login).toHaveBeenCalledWith(expect.any(Object), { username, password });

    expect($router.push).toHaveBeenCalledWith({ name: ROUTES_NAMES.home });
  });

  it('Error shown after trigger submit button with error', async () => {
    login.mockRejectedValueOnce({});

    const wrapper = factory({
      store: createMockedStoreModules([authModule, infoModule]),
      mocks: {
        $route: { query: {} },
        $router: {},
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(login).toHaveBeenCalledWith(expect.any(Object), { username: '', password: '' });

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

    expect(login).toHaveBeenCalledWith(expect.any(Object), {
      username: '',
      password: '',
    });

    expect($router.push).toHaveBeenCalledWith(redirectUrl);
  });

  it('Form submitted after trigger submit button with default view', async () => {
    const defaultView = {
      _id: Faker.datatype.string(),
    };
    currentUser.mockReturnValueOnce({
      defaultview: defaultView,
    });
    const wrapper = factory({
      store: createMockedStoreModules([authModule, infoModule]),
      mocks: {
        $route: { query: {} },
        $router: {},
      },
    });

    selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(login).toHaveBeenCalledWith(expect.any(Object), {
      username: '',
      password: '',
    });

    expect($router.push).toHaveBeenCalledWith({
      name: ROUTES_NAMES.view,
      params: { id: defaultView._id },
    });
  });
});
