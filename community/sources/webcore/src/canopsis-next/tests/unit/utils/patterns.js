import { generateRenderer, generateShallowRenderer, flushPromises } from '@unit/utils/vue';
import { createMockedStoreModules, createPatternsFieldsModule, createPbehaviorPatternsModule } from '@unit/utils/store';

export const generateEntityPatternsTests = (Component, name, customProps = {}) => {
  const stubs = {
    'c-patterns-field': true,
  };

  const {
    patternsFieldsModule,
    fetchDeclareTicketRulePatternFields,
    fetchFlappingRulePatternFields,
    fetchIdleRulePatternFields,
    fetchLinkRulePatternFields,
    fetchResolveRulePatternFields,
    fetchPbehaviorPatternFields,
    fetchAlarmTagPatternFields,
    fetchWidgetFilterPatternFields,
    fetchServicePatternFields,
    fetchStateSettingPatternFields,
    fetchEventFilterPatternFields,
    fetchScenarioPatternFields,
    fetchMetaalarmrulePatternFields,
    fetchInstructionPatternFields,
    fetchKpiFilterPatternFields,
    fetchDynamicInfosPatternFields,
    fetchEventRecordPatternFields,
  } = createPatternsFieldsModule();

  const mockPatternFieldsResponse = {
    entity_pattern: [
      { name: 'name', enabled: true, alias: false },
      { name: 'category', enabled: true, alias: false },
      { name: 'component', enabled: true, alias: false },
    ],
    alarm_pattern: [
      { name: 'output', enabled: true, alias: false },
      { name: 'component', enabled: true, alias: false },
      { name: 'ack', enabled: true, alias: false },
      { name: 'state', enabled: true, alias: false },
    ],
    event_pattern: [],
    pbehavior_pattern: [],
    weather_service_pattern: [],
  };

  fetchDeclareTicketRulePatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchFlappingRulePatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchIdleRulePatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchLinkRulePatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchResolveRulePatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchPbehaviorPatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchAlarmTagPatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchWidgetFilterPatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchServicePatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchStateSettingPatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchEventFilterPatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchScenarioPatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchMetaalarmrulePatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchInstructionPatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchKpiFilterPatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchDynamicInfosPatternFields.mockResolvedValue(mockPatternFieldsResponse);
  fetchEventRecordPatternFields.mockResolvedValue(mockPatternFieldsResponse);

  const { pbehaviorPatternsModule } = createPbehaviorPatternsModule();

  const store = createMockedStoreModules([
    patternsFieldsModule,
    pbehaviorPatternsModule,
  ]);

  const selectPatternsField = wrapper => wrapper.find('c-patterns-field-stub');

  describe(name, () => {
    const factory = generateShallowRenderer(Component, { stubs, store });
    const snapshotFactory = generateRenderer(Component, { stubs, store });

    test('Patterns changed after trigger patterns field', () => {
      const wrapper = factory();

      const patternsField = selectPatternsField(wrapper);

      const newPatterns = {
        alarm_pattern: {},
        pbehavior_pattern: {},
        entity_pattern: {},
      };

      patternsField.triggerCustomEvent('input', newPatterns);

      expect(wrapper).toEmitInput(newPatterns);
    });

    test(`Renders \`${name}\` with default props`, async () => {
      const wrapper = snapshotFactory();

      await flushPromises();

      expect(wrapper).toMatchSnapshot();
    });

    test(`Renders \`${name}\` with custom props`, async () => {
      const wrapper = snapshotFactory({
        propsData: {
          form: {
            alarm_pattern: {},
            entity_pattern: {},
          },
          ...customProps,
        },
      });

      await flushPromises();

      expect(wrapper).toMatchSnapshot();
    });
  });
};
