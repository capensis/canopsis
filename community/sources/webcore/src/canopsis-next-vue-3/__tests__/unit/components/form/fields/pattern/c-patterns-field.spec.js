import { mockModals } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import {
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  PBEHAVIOR_PATTERN_FIELDS,
  EVENT_FILTER_PATTERN_FIELDS,
  PATTERN_CONDITIONS,
  PATTERN_OPTIMIZATION_STATUSES,
  PATTERNS_FIELDS,
} from '@/constants';

import { filterPatternsToForm } from '@/helpers/entities/filter/form';
import { patternToForm } from '@/helpers/entities/pattern/form';

import CPatternsField from '@/components/forms/fields/pattern/c-patterns-field.vue';

const {
  optimizeEntities,
  fetchOptimizeEntitiesStatus,
  updateOptimization,
  removeOptimization,
  checkPatternsEntitiesCount,
  checkPatternsAlarmsCount,
  fetchContextEntitiesListWithoutStore,
} = vi.hoisted(() => ({
  optimizeEntities: vi.fn().mockResolvedValue({}),
  fetchOptimizeEntitiesStatus: vi.fn().mockResolvedValue({}),
  updateOptimization: vi.fn().mockResolvedValue({}),
  removeOptimization: vi.fn().mockResolvedValue({}),
  checkPatternsEntitiesCount: vi.fn().mockResolvedValue({}),
  checkPatternsAlarmsCount: vi.fn().mockResolvedValue({}),
  fetchContextEntitiesListWithoutStore: vi.fn().mockResolvedValue({ data: [], meta: {} }),
}));

vi.mock('@/hooks/store/modules/pattern-entities-optimize', () => ({
  usePatternEntitiesOptimize: () => ({
    optimizeEntities,
    fetchOptimizeEntitiesStatus,
    updateOptimization,
    removeOptimization,
  }),
}));

vi.mock('@/hooks/store/modules/pattern', () => ({
  usePattern: () => ({
    checkPatternsEntitiesCount,
    checkPatternsAlarmsCount,
  }),
}));

vi.mock('@/hooks/store/modules/entity', () => ({
  useEntity: () => ({
    fetchContextEntitiesListWithoutStore,
  }),
}));

const stubs = {
  'c-label': true,
  'c-collapse-panel': true,
  'c-alarm-patterns-field': true,
  'c-service-weather-patterns-field': true,
  'c-entity-patterns-field': {
    name: 'CEntityPatternsField',
    props: ['patterns', 'attributes', 'entityCounter'],
    emits: ['update:value', 'show:entities'],
    template: '<div />',
  },
  'c-pbehavior-patterns-field': true,
  'c-event-filter-patterns-field': true,
  'c-alert': true,
  'c-progress-overlay': true,
  'pattern-try-optimization': {
    name: 'PatternTryOptimization',
    emits: ['try:optimization'],
    template: '<pattern-try-optimization-stub />',
  },
  'pattern-progress': {
    name: 'PatternProgress',
    props: ['failedReason'],
    emits: ['cancel', 'close', 'try-again'],
    template: '<pattern-progress-stub />',
  },
  'pattern-suggestions': true,
  'pattern-count-message': {
    name: 'PatternCountMessage',
    props: ['errorMessage', 'alarmCounter', 'entityCounter'],
    template: '<pattern-count-message-stub />',
  },
  'pattern-field-suggestions-wrapper': {
    name: 'PatternFieldSuggestionsWrapper',
    props: [
      'suggestions',
      'patterns',
      'entityAttributes',
      'optimizedFieldsRegexps',
      'entitiesCount',
    ],
    emits: ['apply:suggestion', 'reject:all', 'show:entities-comparison'],
    template:
      '<pattern-field-suggestions-wrapper-stub><slot /></pattern-field-suggestions-wrapper-stub>',
  },
};

const selectAlarmPatternsField = wrapper => wrapper.find('c-alarm-patterns-field-stub');
const selectEntityPatternsField = wrapper =>
  wrapper.findComponent({ name: 'CEntityPatternsField' });
const selectPbehaviorPatternsField = wrapper => wrapper.find('c-pbehavior-patterns-field-stub');
const selectEventFilterPatternsField = wrapper =>
  wrapper.find('c-event-filter-patterns-field-stub');
const selectTryOptimization = wrapper => wrapper.findComponent({ name: 'PatternTryOptimization' });
const selectOptimizationProgress = wrapper => wrapper.findComponent({ name: 'PatternProgress' });
const selectPatternFieldSuggestionsWrapper = wrapper =>
  wrapper.findComponent({ name: 'PatternFieldSuggestionsWrapper' });
const selectPatternCountMessage = wrapper => wrapper.findComponent({ name: 'PatternCountMessage' });

const createPatternWithRegexpInfos = () => {
  const entityPattern = patternToForm({
    entity_pattern: [
      [
        {
          field: ENTITY_PATTERN_FIELDS.infos,
          cond: {
            type: PATTERN_CONDITIONS.regexp,
            value: 'test.*pattern',
          },
        },
      ],
    ],
  });

  return {
    ...filterPatternsToForm(),
    entity_pattern: entityPattern,
  };
};

const createOptimizationSuggestion = () => ({
  [PATTERNS_FIELDS.entity]: [
    [
      {
        field: ENTITY_PATTERN_FIELDS.id,
        cond: {
          type: PATTERN_CONDITIONS.equal,
          value: 'optimized-id',
        },
      },
    ],
  ],
});

describe('c-patterns-field', () => {
  const $modals = mockModals();
  const patterns = filterPatternsToForm();
  const factory = generateShallowRenderer(CPatternsField, { stubs });

  afterEach(() => {
    optimizeEntities.mockReset();
    fetchOptimizeEntitiesStatus.mockReset();
    updateOptimization.mockReset();
    removeOptimization.mockReset();
    checkPatternsEntitiesCount.mockReset();
    checkPatternsAlarmsCount.mockReset();
    fetchContextEntitiesListWithoutStore.mockReset();
    optimizeEntities.mockResolvedValue({});
    fetchOptimizeEntitiesStatus.mockResolvedValue({});
    updateOptimization.mockResolvedValue({});
    removeOptimization.mockResolvedValue({});
    checkPatternsEntitiesCount.mockResolvedValue({});
    checkPatternsAlarmsCount.mockResolvedValue({});
    fetchContextEntitiesListWithoutStore.mockResolvedValue({ data: [], meta: {} });
  });

  test('Alarm pattern changed after trigger alarm patterns field', () => {
    const wrapper = factory({
      propsData: {
        value: patterns,
        withAlarm: true,
      },
    });

    const alarmPattern = patternToForm({
      alarm_pattern: [
        [
          {
            field: ALARM_PATTERN_FIELDS.output,
            cond: {
              type: PATTERN_CONDITIONS.notEqual,
              value: 'test',
            },
          },
        ],
      ],
    });

    selectAlarmPatternsField(wrapper).triggerCustomEvent('update:patterns', alarmPattern);

    expect(wrapper).toEmitInput({
      ...patterns,
      alarm_pattern: alarmPattern,
    });
  });

  test('Entity pattern changed after trigger entity patterns field', () => {
    const wrapper = factory({
      propsData: {
        value: patterns,
        withEntity: true,
      },
    });

    const entityPattern = patternToForm({
      entity_pattern: [
        [
          {
            field: ENTITY_PATTERN_FIELDS.id,
            cond: {
              type: PATTERN_CONDITIONS.notEqual,
              value: 'id',
            },
          },
        ],
      ],
    });

    selectEntityPatternsField(wrapper).triggerCustomEvent('update:value', entityPattern);

    expect(wrapper).toEmitInput({
      ...patterns,
      entity_pattern: entityPattern,
    });
  });

  test('Shows the total entity pattern counter after checking the filter', async () => {
    const totalEntityPatternCounter = { count: 50, over_limit: false, ms: 1 };
    const counterMethod = vi.fn().mockResolvedValue({
      entities: {
        total_entity_pattern: totalEntityPatternCounter,
      },
    });
    const totalEntityPattern = patternToForm({
      total_entity_pattern: [
        [
          {
            field: ENTITY_PATTERN_FIELDS.type,
            cond: {
              type: PATTERN_CONDITIONS.equal,
              value: 'resource',
            },
          },
        ],
      ],
    });
    const wrapper = factory({
      propsData: {
        value: {
          ...patterns,
          total_entity_pattern: totalEntityPattern,
        },
        withTotalEntity: true,
        counterMethod,
      },
    });

    await wrapper.vm.checkFilter();
    await flushPromises();

    expect(selectEntityPatternsField(wrapper).props('entityCounter')).toEqual(
      totalEntityPatternCounter,
    );
  });

  test('Fetches entities matching the total entity pattern', async () => {
    const totalEntityPatternRules = [
      [
        {
          field: ENTITY_PATTERN_FIELDS.type,
          cond: {
            type: PATTERN_CONDITIONS.equal,
            value: 'resource',
          },
        },
      ],
    ];
    const totalEntityPattern = patternToForm({
      total_entity_pattern: totalEntityPatternRules,
    });
    const wrapper = factory({
      propsData: {
        value: {
          ...patterns,
          total_entity_pattern: totalEntityPattern,
        },
        withTotalEntity: true,
      },
    });

    selectEntityPatternsField(wrapper).triggerCustomEvent('show:entities');

    const [{ config }] = $modals.show.mock.calls[0];

    await config.fetchList({ page: 1 });

    const [{ params }] = fetchContextEntitiesListWithoutStore.mock.calls[0];

    expect(params.page).toBe(1);
    expect(JSON.parse(params.entity_pattern)).toEqual(totalEntityPatternRules);
  });

  test('Shows the total entity count after checking a newly added entity pattern', async () => {
    const combined = { count: 12, ms: 5 };
    checkPatternsEntitiesCount.mockResolvedValue({ combined });

    const wrapper = factory({
      propsData: {
        value: filterPatternsToForm(),
        withEntity: true,
        entityCountersType: true,
        someRequired: true,
      },
    });

    expect(selectPatternCountMessage(wrapper).props('errorMessage')).toBe(
      'Add at least one pattern',
    );

    await wrapper.setProps({ value: createPatternWithRegexpInfos() });
    await wrapper.vm.checkFilter();
    await flushPromises();

    const countMessage = selectPatternCountMessage(wrapper);

    expect(countMessage.props('errorMessage')).toBeUndefined();
    expect(countMessage.props('entityCounter')).toEqual(combined);
  });

  test('Pbehavior pattern changed after trigger pbehavior patterns field', () => {
    const wrapper = factory({
      propsData: {
        value: patterns,
        withPbehavior: true,
      },
    });

    const pbehaviorPattern = patternToForm({
      pbehavior_pattern: [
        [
          {
            field: PBEHAVIOR_PATTERN_FIELDS.name,
            cond: {
              type: PATTERN_CONDITIONS.equal,
              value: 'name',
            },
          },
        ],
      ],
    });

    selectPbehaviorPatternsField(wrapper).triggerCustomEvent('update:value', pbehaviorPattern);

    expect(wrapper).toEmitInput({
      ...patterns,
      pbehavior_pattern: pbehaviorPattern,
    });
  });

  test('Event filter pattern changed after trigger event filter patterns field', () => {
    const wrapper = factory({
      propsData: {
        value: patterns,
        withEvent: true,
      },
    });

    const eventFilterPattern = patternToForm({
      event_pattern: [
        [
          {
            field: EVENT_FILTER_PATTERN_FIELDS.output,
            cond: {
              type: PATTERN_CONDITIONS.equal,
              value: 'output',
            },
          },
        ],
      ],
    });

    selectEventFilterPatternsField(wrapper).triggerCustomEvent('update:value', eventFilterPattern);

    expect(wrapper).toEmitInput({
      ...patterns,
      event_pattern: eventFilterPattern,
    });
  });

  test('Shows try optimization button when pattern has regexp infos and form has changes', async () => {
    const wrapper = factory({
      propsData: {
        value: filterPatternsToForm(),
        withEntity: true,
      },
    });

    await flushPromises();

    await wrapper.setProps({ value: createPatternWithRegexpInfos() });
    await flushPromises();

    expect(selectTryOptimization(wrapper).exists()).toBe(true);
  });

  test('Triggers optimization when try optimization button is clicked', async () => {
    const optimizationId = 'test-optimization-id';
    optimizeEntities.mockResolvedValue({
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.running,
    });
    fetchOptimizeEntitiesStatus.mockResolvedValue({
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.success,
      suggestions: [],
    });

    const wrapper = factory({
      propsData: {
        value: filterPatternsToForm(),
        withEntity: true,
      },
    });

    await flushPromises();
    await wrapper.setProps({ value: createPatternWithRegexpInfos() });
    await flushPromises();

    selectTryOptimization(wrapper).triggerCustomEvent('try:optimization');
    await flushPromises();

    expect(optimizeEntities).toHaveBeenCalled();
  });

  test('Shows optimization suggestions when optimization succeeds with suggestions', async () => {
    const optimizationId = 'test-optimization-id';
    const suggestions = [createOptimizationSuggestion()];
    optimizeEntities.mockResolvedValue({
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.running,
    });
    fetchOptimizeEntitiesStatus.mockResolvedValue({
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.success,
      suggestions,
      optimized_field_regexps: ['infos'],
    });

    const wrapper = factory({
      propsData: {
        value: filterPatternsToForm(),
        withEntity: true,
      },
    });

    await flushPromises();
    await wrapper.setProps({ value: createPatternWithRegexpInfos() });
    await flushPromises();

    selectTryOptimization(wrapper).triggerCustomEvent('try:optimization');
    await flushPromises();

    const suggestionsWrapper = selectPatternFieldSuggestionsWrapper(wrapper);

    expect(suggestionsWrapper.exists()).toBe(true);
    expect(suggestionsWrapper.props('suggestions')).toHaveLength(1);
  });

  test('Applies suggestion when apply suggestion event is triggered', async () => {
    const optimizationId = 'test-optimization-id';
    const suggestion = createOptimizationSuggestion();
    const patternsWithRegexp = createPatternWithRegexpInfos();
    optimizeEntities.mockResolvedValue({
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.running,
    });
    fetchOptimizeEntitiesStatus.mockResolvedValue({
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.success,
      suggestions: [suggestion],
      optimized_field_regexps: ['infos'],
    });

    const wrapper = factory({
      propsData: {
        value: filterPatternsToForm(),
        withEntity: true,
      },
    });

    await flushPromises();
    await wrapper.setProps({ value: patternsWithRegexp });
    await flushPromises();

    selectTryOptimization(wrapper).triggerCustomEvent('try:optimization');
    await flushPromises();

    selectPatternFieldSuggestionsWrapper(wrapper).triggerCustomEvent('apply:suggestion', 0);
    await flushPromises();

    const emittedValue = wrapper.emitted('update:value')[0][0];

    expect(emittedValue.entity_pattern.groups[0].rules[0]).toEqual(
      expect.objectContaining({
        attribute: ENTITY_PATTERN_FIELDS.id,
        operator: 'equal',
        value: 'optimized-id',
      }),
    );
    expect(emittedValue.alarm_pattern).toEqual(patternsWithRegexp.alarm_pattern);
    expect(emittedValue.pbehavior_pattern).toEqual(patternsWithRegexp.pbehavior_pattern);
  });

  test('Cancels optimization when cancel optimization event is triggered', async () => {
    const optimizationId = 'test-optimization-id';
    optimizeEntities.mockResolvedValue({
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.running,
    });
    fetchOptimizeEntitiesStatus.mockResolvedValue({
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.running,
    });

    const wrapper = factory({
      propsData: {
        value: filterPatternsToForm(),
        withEntity: true,
      },
    });

    await flushPromises();
    await wrapper.setProps({ value: createPatternWithRegexpInfos() });
    await flushPromises();

    selectTryOptimization(wrapper).triggerCustomEvent('try:optimization');
    await flushPromises();

    expect(selectOptimizationProgress(wrapper).exists()).toBe(true);

    selectOptimizationProgress(wrapper).triggerCustomEvent('cancel');
    await flushPromises();

    expect(removeOptimization).toHaveBeenCalledWith({ id: optimizationId });
  });

  test('Expands entity panel when optimization suggestions are received', async () => {
    const optimizationId = 'test-optimization-id';
    optimizeEntities.mockResolvedValue({
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.running,
    });
    fetchOptimizeEntitiesStatus.mockResolvedValue({
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.success,
      suggestions: [createOptimizationSuggestion()],
      optimized_field_regexps: ['infos'],
    });

    const wrapper = factory({
      propsData: {
        value: filterPatternsToForm(),
        withEntity: true,
      },
    });

    await flushPromises();

    expect(wrapper.vm.expanded.entity).toBe(false);

    await wrapper.setProps({ value: createPatternWithRegexpInfos() });
    await flushPromises();

    selectTryOptimization(wrapper).triggerCustomEvent('try:optimization');
    await flushPromises();

    expect(wrapper.vm.expanded.entity).toBe(true);
  });

  test('Does not expand entity panel when optimization suggestions are empty', async () => {
    const optimizationId = 'test-optimization-id';
    optimizeEntities.mockResolvedValue({
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.running,
    });
    fetchOptimizeEntitiesStatus.mockResolvedValue({
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.success,
      suggestions: [],
      optimized_field_regexps: [],
    });

    const wrapper = factory({
      propsData: {
        value: filterPatternsToForm(),
        withEntity: true,
      },
    });

    await flushPromises();

    expect(wrapper.vm.expanded.entity).toBe(false);

    await wrapper.setProps({ value: createPatternWithRegexpInfos() });
    await flushPromises();

    selectTryOptimization(wrapper).triggerCustomEvent('try:optimization');
    await flushPromises();

    expect(wrapper.vm.expanded.entity).toBe(false);
  });
});
