import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

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

    expect(wrapper).toEmitInput({ ...form, duration });
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

    expect(wrapper).toEmitInput({ ...form, comment });
  });
});
