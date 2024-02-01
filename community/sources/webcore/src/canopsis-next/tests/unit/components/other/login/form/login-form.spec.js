import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createNumberInputStub } from '@unit/stubs/input';

import LoginForm from '@/components/other/login/form/login-form.vue';

const stubs = {
  'c-password-field': true,
  'v-text-field': createNumberInputStub('v-text-field'),
};

const selectTextField = wrapper => wrapper.find('input.v-text-field');
const selectPasswordField = wrapper => wrapper.find('c-password-field-stub');

describe('login-form', () => {
  const factory = generateShallowRenderer(LoginForm, { stubs });
  const snapshotFactory = generateRenderer(LoginForm, { stubs });

  it('Username changed after trigger text field', () => {
    const wrapper = factory({
      propsData: {
        username: '',
        password: '',
      },
    });

    const username = Faker.datatype.string();

    selectTextField(wrapper).triggerCustomEvent('input', username);

    expect(wrapper).toEmit('input', username);
  });

  it('Password changed after trigger password field', () => {
    const wrapper = factory({
      propsData: {
        username: '',
        password: '',
      },
    });

    const password = Faker.datatype.string();

    selectPasswordField(wrapper).triggerCustomEvent('input', password);

    expect(wrapper).toEmit('input', password);
  });

  it('Renders `login-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `login-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          username: 'username',
          password: 'password',
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
