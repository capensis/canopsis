import { createButtonStub } from '@unit/stubs/button';
import { createFormStub } from '@unit/stubs/form';
import { mockRouter } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import { ROUTES_NAMES } from '@/constants';

import BasicLogin from '@/components/other/login/basic-login.vue';

const auth = vi.hoisted(() => ({
  currentUser: { value: {}, __v_isRef: true },
  login: vi.fn().mockResolvedValue(),
}));
const router = vi.hoisted(() => ({
  route: { query: {} },
  push: vi.fn(),
}));
const form = vi.hoisted(() => ({
  values: { value: { username: '', password: '' }, __v_isRef: true },
}));

vi.mock('@/hooks/auth', () => ({
  useAuth: () => auth,
}));

vi.mock('@/hooks/store/modules/info', () => ({
  useInfo: () => ({ isLDAPAuthEnabled: false, footer: '' }),
}));

vi.mock('@/hooks/form', () => ({
  useForm: () => ({
    values: form.values,
    isSubmitting: false,
    handleSubmit: callback => callback,
  }),
}));

vi.mock('vue-router', () => ({
  useRoute: () => router.route,
  useRouter: () => ({ push: router.push }),
}));

const stubs = {
  'ldap-login-information': true,
  'login-form': true,
  'c-compiled-template': true,
  'c-alert': true,
  'v-btn': createButtonStub('v-btn'),
  'v-form': createFormStub('v-form'),
};

const selectSubmitButton = wrapper => wrapper.find('button.v-btn');
const selectAlert = wrapper => wrapper.find('c-alert-stub');

describe('basic-login', () => {
  const $router = mockRouter();

  const factory = generateShallowRenderer(BasicLogin, {
    stubs,
    attachTo: document.body,
    mocks: { $router },
  });

  afterEach(() => {
    auth.login.mockReset().mockResolvedValue();
    auth.currentUser.value = {};
    form.values.value = { username: '', password: '' };
    router.route.query = {};
    router.push.mockReset();
  });

  it('Form submitted after trigger submit button', async () => {
    const wrapper = factory({
      mocks: { $route: { query: {} } },
    });

    const username = Faker.datatype.string();
    const password = Faker.datatype.string();

    form.values.value = { username, password };

    await selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(auth.login).toHaveBeenCalledWith({ username, password });

    expect(router.push).toHaveBeenCalledWith({ name: ROUTES_NAMES.home });
  });

  it('Error shown after trigger submit button with error', async () => {
    auth.login.mockRejectedValueOnce({});

    const wrapper = factory({
      mocks: {
        $route: { query: {} },
      },
    });

    await selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(auth.login).toHaveBeenCalledWith({ username: '', password: '' });

    expect(selectAlert(wrapper).exists()).toBeTruthy();

    expect(router.push).not.toBeCalled();
  });

  it('Form submitted after trigger submit button with redirect query', async () => {
    const redirectUrl = '/redirect/url';
    router.route.query = { redirect: redirectUrl };
    const wrapper = factory({
      mocks: {
        $route: { query: { redirect: redirectUrl } },
      },
    });

    await selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(auth.login).toHaveBeenCalledWith({
      username: '',
      password: '',
    });

    expect(router.push).toHaveBeenCalledWith(redirectUrl);
  });

  it('Form submitted after trigger submit button with default view', async () => {
    const defaultView = {
      _id: Faker.datatype.string(),
    };
    auth.currentUser.value = {
      defaultview: defaultView,
    };
    const wrapper = factory({
      mocks: {
        $route: { query: {} },
      },
    });

    await selectSubmitButton(wrapper).trigger('click');

    await flushPromises();

    expect(auth.login).toHaveBeenCalledWith({
      username: '',
      password: '',
    });

    expect(router.push).toHaveBeenCalledWith({
      name: ROUTES_NAMES.view,
      params: { id: defaultView._id },
    });
  });
});
