import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import LoginForm from '@/components/other/login/login-form.vue';

const stubs = {
  'c-text-field': true,
  'c-password-field': true,
};

const selectTextField = wrapper => wrapper.find('c-text-field-stub');
const selectPasswordField = wrapper => wrapper.find('c-password-field-stub');

describe('login-form', () => {
  const factory = generateShallowRenderer(LoginForm, { stubs });

  it('Username changed after trigger text field', async () => {
    const wrapper = factory({
      propsData: {
        form: {
          username: '',
          password: '',
        },
      },
    });

    const username = Faker.datatype.string();

    await selectTextField(wrapper).triggerCustomEvent('update:modelValue', username);

    expect(globalThis.__TEST_FORM__.getValue('username')).toBe(username);
  });

  it('Password changed after trigger password field', async () => {
    const wrapper = factory({
      propsData: {
        form: {
          username: '',
          password: '',
        },
      },
    });

    const password = Faker.datatype.string();

    await selectPasswordField(wrapper).triggerCustomEvent('update:modelValue', password);

    expect(globalThis.__TEST_FORM__.getValue('password')).toBe(password);
  });
});
