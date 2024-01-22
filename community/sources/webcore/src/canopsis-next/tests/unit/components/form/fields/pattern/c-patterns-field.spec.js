import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createMockedStoreModules, createPatternModule } from '@unit/utils/store';

import {
  ALARM_PATTERN_FIELDS,
  ENTITY_PATTERN_FIELDS,
  PBEHAVIOR_PATTERN_FIELDS,
  EVENT_FILTER_PATTERN_FIELDS,
  PATTERN_CONDITIONS,
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
};

const selectAlarmPatternsField = wrapper => wrapper.find('c-alarm-patterns-field-stub');
const selectEntityPatternsField = wrapper => wrapper.find('c-entity-patterns-field-stub');
const selectPbehaviorPatternsField = wrapper => wrapper.find('c-pbehavior-patterns-field-stub');
const selectEventFilterPatternsField = wrapper => wrapper.find('c-event-filter-patterns-field-stub');

describe('c-patterns-field', () => {
  const patterns = filterPatternsToForm();
  const { patternModule } = createPatternModule();
  const store = createMockedStoreModules([patternModule]);

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

    expect(wrapper).toEmit('input', {
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

    expect(wrapper).toEmit('input', {
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

    expect(wrapper).toEmit('input', {
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

    expect(wrapper).toEmit('input', {
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
});
