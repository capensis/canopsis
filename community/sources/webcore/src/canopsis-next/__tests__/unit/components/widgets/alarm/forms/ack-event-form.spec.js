import { createActivatorElementStub } from '@unit/stubs/vuetify';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import AckEventForm from '@/components/widgets/alarm/forms/ack-event-form.vue';

const stubs = {
  'c-description-field': true,
  'c-checkbox': true,
  'c-tooltip': createActivatorElementStub('c-tooltip'),
};

const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');
const selectCheckboxField = wrapper => wrapper.find('c-checkbox-stub');

describe('ack-event-form', () => {
  const factory = generateShallowRenderer(AckEventForm, {
    stubs,
    attachTo: document.body,
  });

  test('Comment changed after trigger description field', () => {
    const form = {
      output: Faker.datatype.string(),
      ack_resources: Faker.datatype.boolean(),
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const comment = Faker.datatype.string();

    const descriptionField = selectDescriptionField(wrapper);

    descriptionField.triggerCustomEvent('input', comment);

    expect(wrapper).toEmitInput({ ...form, comment });
  });

  test('Ack resources changed after trigger checkbox field', async () => {
    const form = {
      output: Faker.datatype.string(),
      ack_resources: Faker.datatype.boolean(),
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    await flushPromises();

    const ackResources = !form.ack_resources;

    selectCheckboxField(wrapper).triggerCustomEvent('update:modelValue', ackResources);

    expect(globalThis.__TEST_FORM__.values.value).toEqual({
      ...form,
      ack_resources: ackResources,
    });
  });
});
