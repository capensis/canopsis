import Faker from 'faker';

import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import ServicePauseEventForm from '@/components/widgets/service-weather/service-pause-event-form.vue';

const localVue = createVueInstance();

const stubs = {
  'c-pbehavior-reason-field': true,
  'c-description-field': true,
};

const snapshotStubs = {
  'c-pbehavior-reason-field': true,
  'c-description-field': true,
};

const factory = (options = {}) => shallowMount(ServicePauseEventForm, {
  localVue,
  stubs,

  ...options,
});

const snapshotFactory = (options = {}) => mount(ServicePauseEventForm, {
  localVue,
  stubs: snapshotStubs,

  ...options,
});

const selectPbehaviorReasonField = wrapper => wrapper.find('c-pbehavior-reason-field-stub');
const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');

describe('service-pause-event-form', () => {
  test('Reason changed after trigger reason field', () => {
    const form = {
      comment: Faker.datatype.string(),
      reason: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const reason = Faker.datatype.string();

    const pbehaviorReasonField = selectPbehaviorReasonField(wrapper);

    pbehaviorReasonField.vm.$emit('input', reason);

    expect(wrapper).toEmit('input', { ...form, reason });
  });

  test('Comment changed after trigger textarea', () => {
    const form = {
      comment: Faker.datatype.string(),
      reason: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const comment = Faker.datatype.string();

    const descriptionField = selectDescriptionField(wrapper);

    descriptionField.vm.$emit('input', comment);

    expect(wrapper).toEmit('input', { ...form, comment });
  });

  test('Renders `service-pause-event-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `service-pause-event-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          comment: 'comment',
          reason: 'reason',
        },
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
