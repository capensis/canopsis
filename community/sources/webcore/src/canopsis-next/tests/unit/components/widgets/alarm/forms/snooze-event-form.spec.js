import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { TIME_UNITS } from '@/constants';

import SnoozeEventForm from '@/components/widgets/alarm/forms/snooze-event-form.vue';

const stubs = {
  'c-duration-field': true,
  'c-description-field': true,
};

const selectDurationField = wrapper => wrapper.find('c-duration-field-stub');
const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');

describe('snooze-event-form', () => {
  const factory = generateShallowRenderer(SnoozeEventForm, { stubs });
  const snapshotFactory = generateRenderer(SnoozeEventForm, { stubs });

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

    durationField.triggerCustomEvent('input', duration);

    expect(wrapper).toEmit('input', { ...form, duration });
  });

  test('Output changed after trigger description field', () => {
    const form = {
      duration: {
        unit: TIME_UNITS.day,
        value: Faker.datatype.number(),
      },
      comment: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: { form },
    });

    const comment = Faker.datatype.string();

    const descriptionField = selectDescriptionField(wrapper);

    descriptionField.triggerCustomEvent('input', comment);

    expect(wrapper).toEmit('input', { ...form, comment });
  });

  test('Renders `snooze-event-form` with required props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
