import Faker from 'faker';
import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createNumberInputStub, createSelectInputStub } from '@unit/stubs/input';
import { TIME_UNITS } from '@/constants';
import CDurationField from '@/components/forms/fields/c-duration-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-text-field': createNumberInputStub('v-text-field'),
  'v-select': createSelectInputStub('v-select'),
};

const factory = (options = {}) => shallowMount(CDurationField, {
  localVue,
  stubs,
  ...options,
});

const snapshotFactory = (options = {}) => mount(CDurationField, {
  localVue,
  parentComponent: {
    $_veeValidate: {
      validator: 'new',
    },
  },
  ...options,
});

describe('c-duration-field', () => {
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

    const valueElement = wrapper.find('input.v-text-field');

    valueElement.setValue(newValue);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData.value).toBe(newValue);
    expect(eventData.unit).toBe(duration.unit);
  });

  it('Unit changed after trigger select field', () => {
    const duration = {
      value: Faker.datatype.number(),
      unit: TIME_UNITS.week,
    };
    const wrapper = factory({
      propsData: {
        duration,
      },
    });

    const valueElement = wrapper.find('select.v-select');

    valueElement.setValue(TIME_UNITS.month);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData.unit).toBe(TIME_UNITS.month);
    expect(eventData.value).toBe(duration.value);
  });

  it('Renders `c-duration-field` with default props correctly', () => {
    const wrapper = snapshotFactory();

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-duration-field` with custom props correctly', () => {
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

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-duration-field` with errors correctly', async () => {
    const name = 'customName';

    const wrapper = snapshotFactory({
      propsData: {
        value: {
          value: 10,
          unit: TIME_UNITS.week,
        },
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

    const menuContent = wrapper.find('.v-menu__content');

    expect(wrapper.element).toMatchSnapshot();
    expect(menuContent.element).toMatchSnapshot();
  });

  it('Renders `c-duration-field` with duration value as string correctly', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        duration: {
          value: 'asd',
          unit: TIME_UNITS.week,
        },
      },
    });

    const { $validator: validator } = wrapper.vm;

    await validator.validateAll();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-duration-field` with duration value as undefined correctly', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        duration: {
          value: undefined,
          unit: TIME_UNITS.week,
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-duration-field` with value is greater than the minimum correctly', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        duration: {
          value: 101,
          unit: TIME_UNITS.second,
        },
        min: 100,
      },
    });

    const { $validator: validator } = wrapper.vm;

    await validator.validateAll();

    expect(wrapper.element).toMatchSnapshot();
  });
});
