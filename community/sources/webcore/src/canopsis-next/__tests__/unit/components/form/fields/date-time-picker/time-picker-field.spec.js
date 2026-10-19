import { generateShallowRenderer } from '@unit/utils/vue';

import TimePickerField from '@/components/forms/fields/time-picker/time-picker-field.vue';

const stubs = {
  'v-combobox': {
    name: 'v-combobox',
    props: ['modelValue', 'items', 'menuProps', 'customFilter'],
    emits: ['update:search'],
    template: `
      <input
        :value="modelValue"
        class="v-combobox"
        @input="$emit('update:search', $event.target.value)"
      />
    `,
  },
};

const selectCombobox = wrapper => wrapper.findComponent({ name: 'v-combobox' });

describe('time-picker-field', () => {
  const factory = generateShallowRenderer(TimePickerField, { stubs });

  test('Combobox filter return true with exist text', async () => {
    const wrapper = factory();

    expect(wrapper.vm.filter({}, '12:', '12:13')).toBeTruthy();
  });

  test('Combobox filter return false with excluded text', async () => {
    const wrapper = factory();

    expect(wrapper.vm.filter({}, '13:', '12:13')).toBeFalsy();
  });

  test('Value changed after trigger change with empty value on the combobox', async () => {
    const wrapper = factory({
      propsData: {
        value: '12:00',
      },
    });

    wrapper.vm.change('');

    expect(wrapper).not.toHaveBeenEmit('input');
    expect(wrapper).not.toHaveBeenEmit('update:value');
  });

  test('List scrolled correctly without value', async () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const combobox = selectCombobox(wrapper);

    expect(combobox.props('menuProps')).toEqual({
      minWidth: 90,
      maxHeight: 200,
    });
  });

  test('List scrolled correctly with value', async () => {
    const wrapper = factory({
      propsData: {
        value: '12:00',
      },
    });

    const combobox = selectCombobox(wrapper);

    expect(combobox.props('menuProps')).toEqual({
      minWidth: 90,
      maxHeight: 200,
    });
  });

  test('Value changed after trigger change with valid value on the combobox', async () => {
    const wrapper = factory({
      propsData: {
        value: '12:00',
      },
    });

    const newValue = '13:00';

    wrapper.vm.change(newValue);

    expect(wrapper).toEmitInput(newValue);
  });

  test('Value changed after trigger change with round hours on the combobox', async () => {
    const wrapper = factory({
      propsData: {
        value: '12:00',
        roundHours: true,
      },
    });

    const newValue = '13:13';

    wrapper.vm.change(newValue);

    expect(wrapper).toEmitInput('13:00');
  });
});
