import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { getFormGeneralPatternsTabsStub } from '@unit/stubs/form';

import { EVENT_FILTER_TYPES, PATTERN_CUSTOM_ITEM_VALUE } from '@/constants';

import EventFilterForm from '@/components/other/event-filter/form/event-filter-form.vue';

const stubs = {
  'c-enabled-field': true,
  'c-form-general-patterns-tabs': getFormGeneralPatternsTabsStub(),
  'event-filter-general-form': true,
  'event-filter-patterns-form': true,
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectGeneralForm = wrapper => wrapper.find('event-filter-general-form-stub');
const selectPatternsForm = wrapper => wrapper.find('event-filter-patterns-form-stub');

describe('event-filter-form', () => {
  const form = {
    _id: 'event-filter-id',
    type: EVENT_FILTER_TYPES.drop,
    description: 'event-filter-description',
    priority: 2,
    enabled: true,
    patterns: {
      alarm_pattern: {
        id: PATTERN_CUSTOM_ITEM_VALUE,
        groups: [],
      },
    },
    config: {},
    external_data: [],
  };

  const factory = generateShallowRenderer(EventFilterForm, { stubs });
  const snapshotFactory = generateRenderer(EventFilterForm, { stubs });

  test('General form is rendered in general tab', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    expect(selectGeneralForm(wrapper).exists()).toBe(true);
  });

  test('Enabled changed after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    selectEnabledField(wrapper).triggerCustomEvent('input', false);

    expect(wrapper).toEmitInput({
      ...form,
      enabled: false,
    });
  });

  test('General fields changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const updatedForm = {
      ...form,
      description: Faker.datatype.string(),
    };

    selectGeneralForm(wrapper).triggerCustomEvent('input', updatedForm);

    expect(wrapper).toEmitInput(updatedForm);
  });

  test('Patterns changed after trigger patterns form', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const updatedForm = {
      ...form,
      patterns: {
        alarm_pattern: {},
      },
    };

    selectPatternsForm(wrapper).triggerCustomEvent('input', updatedForm);

    expect(wrapper).toEmitInput(updatedForm);
  });

  test('Renders `event-filter-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `event-filter-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test.each(
    Object.values(EVENT_FILTER_TYPES),
  )('Renders `event-filter-form` with `%s` type', (type) => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          ...form,
          type,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
