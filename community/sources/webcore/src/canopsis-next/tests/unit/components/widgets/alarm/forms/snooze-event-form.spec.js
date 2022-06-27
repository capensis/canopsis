import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';
import { TIME_UNITS } from '@/constants';

import SnoozeEventForm from '@/components/widgets/alarm/forms/snooze-event-form.vue';

const localVue = createVueInstance();

const stubs = {
  'c-duration-field': true,
  'c-description-field': true,
};

const factory = (options = {}) => shallowMount(SnoozeEventForm, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(SnoozeEventForm, {
  localVue,
  stubs,

  ...options,
});

const selectDurationField = wrapper => wrapper.find('c-duration-field-stub');
const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');

describe('snooze-event-form', () => {
  test('Duration changed after trigger duration field', () => {
    const form = {
      duration: {
        unit: TIME_UNITS.day,
        value: Faker.datatype.number(),
      },
      output: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: { form },
    });

    const duration = {
      unit: TIME_UNITS.hour,
      value: Faker.datatype.number(),
    };

    const durationField = selectDurationField(wrapper);

    durationField.vm.$emit('input', duration);

    expect(wrapper).toEmit('input', { ...form, duration });
  });

  test('Output changed after trigger description field', () => {
    const form = {
      duration: {
        unit: TIME_UNITS.day,
        value: Faker.datatype.number(),
      },
      output: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: { form },
    });

    const output = Faker.datatype.string();

    const descriptionField = selectDescriptionField(wrapper);

    descriptionField.vm.$emit('input', output);

    expect(wrapper).toEmit('input', { ...form, output });
  });

  test('Renders `snooze-event-form` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `snooze-event-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          duration: {
            unit: TIME_UNITS.day,
            value: 2,
          },
          output: 'output',
        },
        isNoteRequired: true,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
