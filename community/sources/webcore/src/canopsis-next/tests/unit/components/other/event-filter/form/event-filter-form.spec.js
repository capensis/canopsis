import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { EVENT_FILTER_ENRICHMENT_AFTER_TYPES, EVENT_FILTER_TYPES, PATTERN_CUSTOM_ITEM_VALUE } from '@/constants';

import EventFilterForm from '@/components/other/event-filter/form/event-filter-form.vue';

const stubs = {
  'c-id-field': true,
  'c-event-filter-type-field': true,
  'c-description-field': true,
  'c-priority-field': true,
  'c-enabled-field': true,
  'c-patterns-field': true,
  'c-information-block': true,
  'c-collapse-panel': true,
  'external-data-form': true,
  'event-filter-change-entity-form': true,
  'event-filter-enrichment-form': true,
  'pbehavior-recurrence-rule-field': true,
  'event-filter-drop-intervals-field': true,
};

const selectIdField = wrapper => wrapper.find('c-id-field-stub');
const selectEventFilterTypeField = wrapper => wrapper.find('c-event-filter-type-field-stub');
const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');
const selectPriorityField = wrapper => wrapper.find('c-priority-field-stub');
const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectPatternsField = wrapper => wrapper.find('c-patterns-field-stub');
const selectEventFilterChangeEntityForm = wrapper => wrapper.find('event-filter-change-entity-form-stub');
const selectEventFilterEnrichmentForm = wrapper => wrapper.find('event-filter-enrichment-form-stub');
const selectEventFilterDropIntervalsField = wrapper => wrapper.find('event-filter-drop-intervals-field-stub');

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

  test('ID changed after trigger id field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const idField = selectIdField(wrapper);

    const newId = Faker.datatype.string();

    idField.triggerCustomEvent('input', newId);

    expect(wrapper).toEmit('input', {
      ...form,
      _id: newId,
    });
  });

  test('Type changed after trigger type field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const eventFilterTypeField = selectEventFilterTypeField(wrapper);

    eventFilterTypeField.triggerCustomEvent('input', EVENT_FILTER_TYPES.enrichment);

    expect(wrapper).toEmit('input', {
      ...form,
      type: EVENT_FILTER_TYPES.enrichment,
    });
  });

  test('Description changed after trigger description field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const descriptionField = selectDescriptionField(wrapper);

    const description = Faker.datatype.string();

    descriptionField.triggerCustomEvent('input', description);

    expect(wrapper).toEmit('input', {
      ...form,
      description,
    });
  });

  test('Priority changed after trigger priority field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const priorityField = selectPriorityField(wrapper);

    const priority = Faker.datatype.string();

    priorityField.triggerCustomEvent('input', priority);

    expect(wrapper).toEmit('input', {
      ...form,
      priority,
    });
  });

  test('Enabled changed after trigger enabled field', () => {
    const enabled = Faker.datatype.boolean();
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          enabled,
        },
      },
    });

    const enabledField = selectEnabledField(wrapper);

    const newEnabled = !enabled;

    enabledField.triggerCustomEvent('input', newEnabled);

    expect(wrapper).toEmit('input', {
      ...form,
      enabled: newEnabled,
    });
  });

  test('Patterns changed after trigger patterns field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const patternsField = selectPatternsField(wrapper);

    const newPatterns = {
      id: Faker.datatype.string(),
    };

    patternsField.triggerCustomEvent('input', newPatterns);

    expect(wrapper).toEmit('input', {
      ...form,
      patterns: newPatterns,
    });
  });

  test('Change entity config changed after trigger change entity form', () => {
    const changeEntityForm = {
      ...form,
      type: EVENT_FILTER_TYPES.changeEntity,
    };
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          type: EVENT_FILTER_TYPES.changeEntity,
        },
      },
    });

    const eventFilterChangeEntityForm = selectEventFilterChangeEntityForm(wrapper);

    const newConfig = {
      resource: Faker.datatype.string(),
      component: Faker.datatype.string(),
      connector: Faker.datatype.string(),
      connector_name: Faker.datatype.string(),
    };

    eventFilterChangeEntityForm.triggerCustomEvent('input', newConfig);

    expect(wrapper).toEmit('input', {
      ...changeEntityForm,
      config: newConfig,
    });
  });

  test('Enrichment config changed after trigger enrichment form', () => {
    const enrichmentForm = {
      ...form,
      type: EVENT_FILTER_TYPES.enrichment,
    };
    const wrapper = factory({
      propsData: {
        form: enrichmentForm,
      },
    });

    const eventFilterEnrichmentForm = selectEventFilterEnrichmentForm(wrapper);

    const updatedForm = {
      ...form,
      config: {
        on_success: EVENT_FILTER_ENRICHMENT_AFTER_TYPES.break,
        on_failure: EVENT_FILTER_ENRICHMENT_AFTER_TYPES.drop,
        actions: [],
      },
      external_data: [],
    };

    eventFilterEnrichmentForm.triggerCustomEvent('input', updatedForm);

    expect(wrapper).toEmit('input', updatedForm);
  });

  test('Drop intervals fields changed after trigger', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const eventFilterDropIntervalsField = selectEventFilterDropIntervalsField(wrapper);

    const newData = {
      ...form,
      exdates: [{}],
      exceptions: [],
    };

    eventFilterDropIntervalsField.triggerCustomEvent('input', newData);

    expect(wrapper).toEmit('input', newData);
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
