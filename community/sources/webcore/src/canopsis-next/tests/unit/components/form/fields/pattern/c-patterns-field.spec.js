import { flushPromises, generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import {
  createMockedStoreModules,
  createPatternModule,
  createAlarmModule,
  createEntityModule,
  createPatternEntitiesOptimizeModule,
} from '@unit/utils/store';

import {
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  PBEHAVIOR_PATTERN_FIELDS,
  EVENT_FILTER_PATTERN_FIELDS,
  PATTERN_CONDITIONS,
  PATTERN_OPTIMIZATION_STATUSES,
} from '@/constants';

import { patternToForm } from '@/helpers/entities/pattern/form';
import { filterPatternsToForm } from '@/helpers/entities/filter/form';

import CPatternsField from '@/components/forms/fields/pattern/c-patterns-field.vue';

const stubs = {
  'c-collapse-panel': true,
  'c-alarm-patterns-field': true,
  'c-service-weather-patterns-field': true,
  'c-entity-patterns-field': true,
  'c-pbehavior-patterns-field': true,
  'c-event-filter-patterns-field': true,
  'c-alert': true,
  'pattern-try-optimization': true,
  'pattern-optimization-progress': true,
  'pattern-suggestions': true,
};

const selectAlarmPatternsField = wrapper => wrapper.find('c-alarm-patterns-field-stub');
const selectEntityPatternsField = wrapper => wrapper.find('c-entity-patterns-field-stub');
const selectPbehaviorPatternsField = wrapper => wrapper.find('c-pbehavior-patterns-field-stub');
const selectEventFilterPatternsField = wrapper => wrapper.find('c-event-filter-patterns-field-stub');
const selectTryOptimization = wrapper => wrapper.find('pattern-try-optimization-stub');
const selectOptimizationProgress = wrapper => wrapper.find('pattern-optimization-progress-stub');
const selectPatternSuggestions = wrapper => wrapper.find('pattern-suggestions-stub');

const createPatternWithRegexpInfos = () => {
  const entityPattern = patternToForm({
    entity_pattern: [[
      {
        field: ENTITY_PATTERN_FIELDS.infos,
        cond: {
          type: PATTERN_CONDITIONS.regexp,
          value: 'test.*pattern',
        },
      },
    ]],
  });

  return {
    ...filterPatternsToForm(),
    entity_pattern: entityPattern,
  };
};

describe('c-patterns-field', () => {
  const patterns = filterPatternsToForm();
  const { patternModule } = createPatternModule();
  const { alarmModule } = createAlarmModule();
  const { entityModule } = createEntityModule();
  const { patternEntitiesOptimizeModule } = createPatternEntitiesOptimizeModule();
  const store = createMockedStoreModules([
    patternModule,
    alarmModule,
    entityModule,
    patternEntitiesOptimizeModule,
  ]);

  const factory = generateShallowRenderer(CPatternsField, { stubs });
  const snapshotFactory = generateRenderer(CPatternsField, { stubs });

  test('Alarm pattern changed after trigger alarm patterns field', () => {
    const wrapper = factory({
      propsData: {
        value: patterns,
        withAlarm: true,
      },
      store,
    });

    const alarmPattern = patternToForm({
      alarm_pattern: [[
        {
          field: ALARM_PATTERN_FIELDS.output,
          cond: {
            type: PATTERN_CONDITIONS.notEqual,
            value: 'test',
          },
        },
      ]],
    });

    const alarmPatternsField = selectAlarmPatternsField(wrapper);

    alarmPatternsField.triggerCustomEvent('input', alarmPattern);

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
      store,
    });

    const entityPattern = patternToForm({
      entity_pattern: [[
        {
          field: ENTITY_PATTERN_FIELDS.id,
          cond: {
            type: PATTERN_CONDITIONS.notEqual,
            value: 'id',
          },
        },
      ]],
    });

    const entityPatternsField = selectEntityPatternsField(wrapper);

    entityPatternsField.triggerCustomEvent('input', entityPattern);

    expect(wrapper).toEmitInput({
      ...patterns,
      entity_pattern: entityPattern,
    });
  });

  test('Pbehavior pattern changed after trigger pbehavior patterns field', () => {
    const wrapper = factory({
      propsData: {
        value: patterns,
        withPbehavior: true,
      },
      store,
    });

    const pbehaviorPattern = patternToForm({
      entity_pattern: [[
        {
          field: PBEHAVIOR_PATTERN_FIELDS.name,
          cond: {
            type: PATTERN_CONDITIONS.equal,
            value: 'name',
          },
        },
      ]],
    });

    const pbehaviorPatternsField = selectPbehaviorPatternsField(wrapper);

    pbehaviorPatternsField.triggerCustomEvent('input', pbehaviorPattern);

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
      store,
    });

    const eventFilterPattern = patternToForm({
      entity_pattern: [[
        {
          field: EVENT_FILTER_PATTERN_FIELDS.output,
          cond: {
            type: PATTERN_CONDITIONS.equal,
            value: 'output',
          },
        },
      ]],
    });

    const eventFilterPatternsField = selectEventFilterPatternsField(wrapper);

    eventFilterPatternsField.triggerCustomEvent('input', eventFilterPattern);

    expect(wrapper).toEmitInput({
      ...patterns,
      event_pattern: eventFilterPattern,
    });
  });

  test('Renders `c-patterns-field` with default props', () => {
    const wrapper = snapshotFactory({ store });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `c-patterns-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: patterns,
        disabled: true,
        withAlarm: true,
        withEvent: true,
        withEntity: true,
        withPbehavior: true,
        withTotalEntity: true,
        withServiceWeather: true,
        required: true,
        readonly: true,
        someRequired: true,
        name: 'name',
      },
      store,
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Shows try optimization button when pattern has regexp infos and form has changes', async () => {
    const initialPatterns = filterPatternsToForm();
    const wrapper = factory({
      propsData: {
        value: initialPatterns,
        withEntity: true,
      },
      store,
    });

    await flushPromises();

    const patternsWithRegexp = createPatternWithRegexpInfos();
    wrapper.setProps({ value: patternsWithRegexp });

    await flushPromises();

    const tryOptimizationComponent = selectTryOptimization(wrapper);

    expect(tryOptimizationComponent.exists()).toBe(true);
  });

  test('Triggers optimization when try optimization button is clicked', async () => {
    const {
      patternEntitiesOptimizeModule: optimizeModule,
      optimize,
      fetchOptimizeStatus,
    } = createPatternEntitiesOptimizeModule();
    const testStore = createMockedStoreModules([
      patternModule,
      alarmModule,
      entityModule,
      optimizeModule,
    ]);

    const optimizationId = 'test-optimization-id';
    const optimizationResponse = {
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.running,
    };
    const successResponse = {
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.success,
      suggestions: [],
    };

    optimize.mockResolvedValue(optimizationResponse);
    fetchOptimizeStatus.mockResolvedValue(successResponse);

    const initialPatterns = filterPatternsToForm();
    const wrapper = factory({
      propsData: {
        value: initialPatterns,
        withEntity: true,
      },
      store: testStore,
    });

    await flushPromises();

    const patternsWithRegexp = createPatternWithRegexpInfos();
    wrapper.setProps({ value: patternsWithRegexp });

    await flushPromises();

    const tryOptimizationComponent = selectTryOptimization(wrapper);

    expect(tryOptimizationComponent.exists()).toBe(true);
    tryOptimizationComponent.triggerCustomEvent('try:optimization');

    await flushPromises();

    expect(optimize).toHaveBeenCalled();
  });

  test('Shows optimization suggestions when optimization succeeds with suggestions', async () => {
    const {
      patternEntitiesOptimizeModule: optimizeModule,
      optimize,
      fetchOptimizeStatus,
    } = createPatternEntitiesOptimizeModule();
    const testStore = createMockedStoreModules([
      patternModule,
      alarmModule,
      entityModule,
      optimizeModule,
    ]);

    const optimizationId = 'test-optimization-id';
    const optimizationResponse = {
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.running,
    };
    const successResponse = {
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.success,
      suggestions: [
        {
          groups: [[
            {
              field: ENTITY_PATTERN_FIELDS.id,
              cond: {
                type: PATTERN_CONDITIONS.equal,
                value: 'optimized-id',
              },
            },
          ]],
        },
      ],
      optimized_field_regexps: ['infos'],
    };

    optimize.mockResolvedValue(optimizationResponse);
    fetchOptimizeStatus.mockResolvedValue(successResponse);

    const initialPatterns = filterPatternsToForm();
    const wrapper = factory({
      propsData: {
        value: initialPatterns,
        withEntity: true,
      },
      store: testStore,
    });

    await flushPromises();

    const patternsWithRegexp = createPatternWithRegexpInfos();
    wrapper.setProps({ value: patternsWithRegexp });

    await flushPromises();

    const tryOptimizationComponent = selectTryOptimization(wrapper);

    expect(tryOptimizationComponent.exists()).toBe(true);
    tryOptimizationComponent.triggerCustomEvent('try:optimization');

    await flushPromises();

    const suggestionsComponent = selectPatternSuggestions(wrapper);

    expect(suggestionsComponent.exists()).toBe(true);
  });

  test('Applies suggestion when apply suggestion event is triggered', async () => {
    const {
      patternEntitiesOptimizeModule: optimizeModule,
      optimize,
      fetchOptimizeStatus,
    } = createPatternEntitiesOptimizeModule();
    const testStore = createMockedStoreModules([
      patternModule,
      alarmModule,
      entityModule,
      optimizeModule,
    ]);

    const optimizationId = 'test-optimization-id';
    const optimizationResponse = {
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.running,
    };
    const successResponse = {
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.success,
      suggestions: [
        {
          groups: [[
            {
              field: ENTITY_PATTERN_FIELDS.id,
              cond: {
                type: PATTERN_CONDITIONS.equal,
                value: 'optimized-id',
              },
            },
          ]],
        },
      ],
      optimized_field_regexps: ['infos'],
    };

    optimize.mockResolvedValue(optimizationResponse);
    fetchOptimizeStatus.mockResolvedValue(successResponse);

    const initialPatterns = filterPatternsToForm();
    const wrapper = factory({
      propsData: {
        value: initialPatterns,
        withEntity: true,
      },
      store: testStore,
    });

    await flushPromises();

    const patternsWithRegexp = createPatternWithRegexpInfos();
    wrapper.setProps({ value: patternsWithRegexp });

    await flushPromises();

    const tryOptimizationComponent = selectTryOptimization(wrapper);

    expect(tryOptimizationComponent.exists()).toBe(true);
    tryOptimizationComponent.triggerCustomEvent('try:optimization');

    await flushPromises();

    const suggestionsComponent = selectPatternSuggestions(wrapper);

    expect(suggestionsComponent.exists()).toBe(true);
    suggestionsComponent.triggerCustomEvent('apply:suggestion', 0);

    await flushPromises();

    expect(wrapper).toEmitInput();
  });

  test('Cancels optimization when cancel optimization event is triggered', async () => {
    const {
      patternEntitiesOptimizeModule: optimizeModule,
      optimize,
      fetchOptimizeStatus,
      remove,
    } = createPatternEntitiesOptimizeModule();
    const testStore = createMockedStoreModules([
      patternModule,
      alarmModule,
      entityModule,
      optimizeModule,
    ]);

    const optimizationId = 'test-optimization-id';
    const optimizationResponse = {
      _id: optimizationId,
      status: PATTERN_OPTIMIZATION_STATUSES.running,
    };

    optimize.mockResolvedValue(optimizationResponse);
    fetchOptimizeStatus.mockResolvedValue(optimizationResponse);
    remove.mockResolvedValue({});

    const initialPatterns = filterPatternsToForm();
    const wrapper = factory({
      propsData: {
        value: initialPatterns,
        withEntity: true,
      },
      store: testStore,
    });

    await flushPromises();

    const patternsWithRegexp = createPatternWithRegexpInfos();
    wrapper.setProps({ value: patternsWithRegexp });

    await flushPromises();

    const tryOptimizationComponent = selectTryOptimization(wrapper);

    expect(tryOptimizationComponent.exists()).toBe(true);
    tryOptimizationComponent.triggerCustomEvent('try:optimization');

    await flushPromises();

    const progressComponent = selectOptimizationProgress(wrapper);

    expect(progressComponent.exists()).toBe(true);

    progressComponent.triggerCustomEvent('cancel:optimization');

    await flushPromises();

    expect(remove).toHaveBeenCalledWith(expect.any(Object), { id: optimizationId });
  });
});
