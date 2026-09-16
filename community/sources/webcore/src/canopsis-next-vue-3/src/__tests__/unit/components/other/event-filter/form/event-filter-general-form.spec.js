import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { EVENT_FILTER_TYPES } from '@/constants';

import EventFilterGeneralForm from '@/components/other/event-filter/form/event-filter-general-form.vue';

const stubs = {
  'c-name-field': true,
  'c-form-block': true,
  'c-form-block-row': true,
  'c-event-filter-type-field': true,
  'c-priority-field': true,
  'c-description-field': true,
  'event-filter-drop-intervals-field': true,
  'pbehavior-recurrence-rule-field': true,
  'event-filter-enrichment-form': true,
  'event-filter-change-entity-form': true,
};

const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');
const selectEnrichmentForm = wrapper => wrapper.find('event-filter-enrichment-form-stub');
const selectChangeEntityForm = wrapper => wrapper.find('event-filter-change-entity-form-stub');

describe('event-filter-general-form', () => {
  const form = {
    type: EVENT_FILTER_TYPES.drop,
    priority: 2,
    description: Faker.datatype.string(),
    config: {},
  };

  const factory = generateShallowRenderer(EventFilterGeneralForm, { stubs });

  test('Description changed after trigger description field', () => {
    const wrapper = factory({
      propsData: {
        form,
      },
    });

    const newDescription = Faker.datatype.string();

    selectDescriptionField(wrapper).triggerCustomEvent('input', newDescription);

    expect(wrapper).toEmitInput({ ...form, description: newDescription });
  });

  test('Enrichment form is rendered for enrichment type', () => {
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          type: EVENT_FILTER_TYPES.enrichment,
        },
      },
    });

    expect(selectEnrichmentForm(wrapper).exists()).toBe(true);
    expect(selectChangeEntityForm(wrapper).exists()).toBe(false);
  });

  test('Change entity form is rendered for change entity type', () => {
    const wrapper = factory({
      propsData: {
        form: {
          ...form,
          type: EVENT_FILTER_TYPES.changeEntity,
        },
      },
    });

    expect(selectChangeEntityForm(wrapper).exists()).toBe(true);
    expect(selectEnrichmentForm(wrapper).exists()).toBe(false);
  });
});
