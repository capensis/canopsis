import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createInputStub } from '@unit/stubs/input';
import { COLORS } from '@/config';
import { DATETIME_FORMATS } from '@/constants';

import CDatePickerField from '@/components/forms/fields/date-picker/c-date-picker-field.vue';

const stubs = {
  'v-menu': {
    template: '<div><slot name="activator"/><slot/></div>',
  },
  'v-date-picker': createInputStub('v-date-picker'),
  'v-text-field': createInputStub('v-text-field'),
};

const selectDatePicker = wrapper => wrapper.find('.v-date-picker');
const selectTextField = wrapper => wrapper.find('.v-text-field');

describe('c-date-picker-field', () => {
  const factory = generateShallowRenderer(CDatePickerField, {

    stubs,
    attachTo: document.body,
  });
  const snapshotFactory = generateRenderer(CDatePickerField, {

    attachTo: document.body,
  });

  test('Value changed after trigger date picker', () => {
    const wrapper = factory({
      propsData: {
        value: '2022-10-03',
      },
    });

    const newValue = '2022-10-04';

    selectDatePicker(wrapper).vm.$emit('input', newValue);

    expect(wrapper).toEmit('input', newValue);
  });

  test('Value cleared after trigger text field', () => {
    const wrapper = factory({
      propsData: {
        value: '2022-10-03',
      },
    });

    selectTextField(wrapper).vm.$emit('click:append');

    expect(wrapper).toEmit('input', null);
  });

  test('Change event emitted after trigger date picker', () => {
    const wrapper = factory({
      propsData: {
        value: '',
      },
    });

    const newValue = '2022-10-04';

    selectDatePicker(wrapper).vm.$emit('change', newValue);

    expect(wrapper).toEmit('change', newValue);
  });

  test('Renders `c-date-picker-field` with default props', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `c-date-picker-field` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 123123123,
        label: 'Custom label',
        name: 'customName',
        color: COLORS.secondary,
        format: DATETIME_FORMATS.long,
        error: true,
        hideDetails: true,
        disabled: true,
        min: 12312312,
        max: 1231231230,
        allowedDates: () => true,
        clearable: true,
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `c-date-picker-field` with slots', async () => {
    const wrapper = snapshotFactory({
      slots: {
        append: '<div class="append-slot" />',
      },
    });

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  test('Renders `c-date-picker-field` with errors', async () => {
    const name = 'date-picker';
    const wrapper = snapshotFactory({
      propsData: {
        value: '2022-01-12',
        name,
      },
      parentComponent: {
        $_veeValidate: {
          validator: 'new',
        },
      },
    });

    wrapper.getValidator().errors.add([
      {
        field: name,
        msg: 'Value error',
      },
    ]);

    await flushPromises();

    expect(document.body.innerHTML).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });
});
