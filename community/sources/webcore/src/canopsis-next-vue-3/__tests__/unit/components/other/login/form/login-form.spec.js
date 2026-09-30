import { createNumberInputStub } from '@unit/stubs/input';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import LoginForm from '@/components/other/login/form/login-form.vue';

const stubs = {
  'c-password-field': true,
  'v-text-field': createNumberInputStub('v-text-field'),
};

const selectTextField = wrapper => wrapper.find('input.v-text-field');
const selectPasswordField = wrapper => wrapper.find('c-password-field-stub');

describe('login-form', () => {
  const factory = generateShallowRenderer(LoginForm, { stubs });

  it('Username changed after trigger text field', () => {
    const wrapper = factory({
      propsData: {
        username: '',
        password: '',
      },
    });

    const username = Faker.datatype.string();

    selectTextField(wrapper).triggerCustomEvent('input', username);

    expect(wrapper).toEmitInput(username);
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

    expect(wrapper).toEmitInput(password);
  });
});
