import { generateShallowRenderer } from '@unit/utils/vue';

import EventFilterEnrichmentActionsField from '@/components/other/event-filter/form/fields/event-filter-enrichment-actions-field.vue';

const stubs = {
  'c-draggable-list-field': {
    props: ['modelValue'],
    template: '<div><slot name="item" :element="modelValue[0]" :index="0" /></div>',
  },
  'event-filter-enrichment-action-field': true,
};

describe('event-filter-enrichment-actions-field', () => {
  const factory = generateShallowRenderer(EventFilterEnrichmentActionsField, { stubs });

  test('uses the action key in field names', () => {
    const wrapper = factory({
      propsData: {
        actions: [{ key: 'stable-action-key' }],
      },
    });

    expect(wrapper.find('event-filter-enrichment-action-field-stub').attributes('name')).toBe(
      'config.actions.stable-action-key',
    );
  });
});
