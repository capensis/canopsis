import Faker from 'faker';
import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createNumberInputStub } from '@unit/stubs/input';
import { TIME_UNITS } from '@/constants';

import CDurationField from '@/components/forms/fields/duration/c-duration-field.vue';

const stubs = {
  'c-number-field': createNumberInputStub('c-number-field'),
  'v-select': {
    props: ['value'],
    template: `
      <input
        :value="value"
        class="v-select"
        @input="$listeners.input($event.target.value)"
      />
    `,
  },
};

const snapshotStubs = {
  'c-number-field': true,
};

const selectNumberField = wrapper => wrapper.find('input.c-number-field');
const selectSelectField = wrapper => wrapper.find('.v-select');

describe('c-duration-field', () => {
  const factory = generateShallowRenderer(CDurationField, {

    stubs,
  });

  const snapshotFactory = generateRenderer(CDurationField, {

    stubs: snapshotStubs,
    parentComponent: {
      $_veeValidate: {
        validator: 'new',
      },
    },
  });

  it('Value changed after trigger text field', () => {
    const duration = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.week,
    };
    const wrapper = factory({
      propsData: {
        duration,
      },
    });
    const newValue = Faker.datatype.number();

    const valueElement = selectNumberField(wrapper);

    valueElement.setValue(newValue);

    expect(wrapper).toEmit('input', {
      unit: duration.unit,
      value: newValue,
    });
  });

  it('Unit changed after trigger select field', () => {
    const duration = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.week,
    };
    const wrapper = factory({
      propsData: {
        duration,
        long: true,
      },
    });

    const valueElement = selectSelectField(wrapper);

    valueElement.vm.$emit('change', TIME_UNITS.month);

    expect(wrapper).toEmit('input', {
      unit: TIME_UNITS.month,
      value: duration.value,
    });
  });

  it('Value cleared after trigger select field without value', () => {
    const duration = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.week,
    };
    const wrapper = factory({
      propsData: {
        duration,
        long: true,
      },
    });

    const valueElement = selectSelectField(wrapper);

    valueElement.vm.$emit('change');

    expect(wrapper).toEmit('input', {
      unit: undefined,
      value: undefined,
    });
  });

  it('Renders `c-duration-field` with default props', () => {
    const wrapper = snapshotFactory();

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-duration-field` with long time list', () => {
    const wrapper = snapshotFactory({
      propsData: {
        long: true,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-duration-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        duration: {
          value: 12,
          unit: TIME_UNITS.week,
        },
        label: 'Custom label',
        unitsLabel: 'Custom units label',
        units: [{
          value: TIME_UNITS.week,
          text: 'Week',
        }, {
          value: TIME_UNITS.month,
          text: 'Month',
        }],
        name: 'customName',
        disabled: true,
        required: true,
        clearable: true,
        min: 10,
      },
    });

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-duration-field` with errors', async () => {
    const name = 'customName';

    const wrapper = snapshotFactory({
      propsData: {
        value: {
          value: 10,
          unit: TIME_UNITS.week,
        },
        long: true,
        name,
      },
    });

    const { $validator: validator } = wrapper.vm;

    validator.errors.add([
      {
        field: `${name}.value`,
        msg: 'Value error',
      },
      {
        field: `${name}.unit`,
        msg: 'Unit error',
      },
    ]);

    await validator.validateAll();

    const menuContent = wrapper.findMenu();

    expect(wrapper).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-duration-field` with duration value as string', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        duration: {
          value: 'asd',
          unit: TIME_UNITS.week,
        },
        long: true,
      },
    });

    const { $validator: validator } = wrapper.vm;

    await validator.validateAll();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-duration-field` with duration value as undefined', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        duration: {
          value: undefined,
          unit: TIME_UNITS.week,
        },
        long: true,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `c-duration-field` with value is greater than the minimum', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        duration: {
          value: 101,
          unit: TIME_UNITS.second,
        },
        min: 100,
        long: true,
      },
    });

    const { $validator: validator } = wrapper.vm;

    await validator.validateAll();

    expect(wrapper).toMatchSnapshot();
  });
});
