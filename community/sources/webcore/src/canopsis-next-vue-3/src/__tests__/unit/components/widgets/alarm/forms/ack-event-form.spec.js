import { createCheckboxInputStub } from '@unit/stubs/input';
import { createActivatorElementStub } from '@unit/stubs/vuetify';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';
import Faker from 'faker';

import AckEventForm from '@/components/widgets/alarm/forms/ack-event-form.vue';

const stubs = {
  'c-description-field': true,
  'v-checkbox': createCheckboxInputStub('v-checkbox'),
  'v-tooltip': createActivatorElementStub('v-tooltip'),
};

const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');
const selectCheckboxField = wrapper => wrapper.find('.v-checkbox');

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

    selectCheckboxField(wrapper).triggerCustomEvent('change', ackResources);

    expect(wrapper).toEmitInput({ ...form, ack_resources: ackResources });
  });
});
