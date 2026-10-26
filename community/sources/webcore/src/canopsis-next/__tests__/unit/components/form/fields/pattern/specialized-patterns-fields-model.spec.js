import { generateShallowRenderer } from '@unit/utils/vue';

import CEntityPatternsField from '@/components/forms/fields/entity/c-entity-patterns-field.vue';
import CEventFilterPatternsField from '@/components/forms/fields/event-filter/c-event-filter-patterns-field.vue';
import CPbehaviorPatternsField from '@/components/forms/fields/pbehavior/c-pbehavior-patterns-field.vue';
import CServiceWeatherPatternsField from '@/components/forms/fields/service-weather/c-service-weather-patterns-field.vue';

vi.mock('@/hooks/store/modules/service', () => ({
  useService: () => ({
    fetchEntityInfosKeysWithoutStore: vi.fn().mockResolvedValue({ data: [] }),
  }),
}));

const patternEditorFieldStub = {
  name: 'PatternEditorField',
  props: ['patterns'],
  emits: ['update:patterns'],
  template: '<div />',
};

const stubs = {
  'pattern-editor-field': patternEditorFieldStub,
};

describe('specialized patterns fields v-model', () => {
  const factories = {
    entity: generateShallowRenderer(CEntityPatternsField, { stubs }),
    event: generateShallowRenderer(CEventFilterPatternsField, { stubs }),
    pbehavior: generateShallowRenderer(CPbehaviorPatternsField, { stubs }),
    serviceWeather: generateShallowRenderer(CServiceWeatherPatternsField, { stubs }),
  };

  test.each(Object.entries(factories))(
    '%s field emits update:patterns',
    (_, factory) => {
      const patterns = { id: 'custom', groups: [] };
      const nextPatterns = { ...patterns, groups: [{ key: 'group', rules: [] }] };
      const wrapper = factory({ propsData: { patterns } });

      wrapper.findComponent(patternEditorFieldStub).vm.$emit('update:patterns', nextPatterns);

      expect(wrapper.emitted('update:patterns')).toEqual([[nextPatterns]]);
      expect(wrapper.emitted('update:value')).toBeUndefined();
    },
  );
});
