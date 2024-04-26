import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { ALARM_STATES, ENTITY_PATTERN_FIELDS, PATTERN_CONDITIONS } from '@/constants';

import { serviceToForm } from '@/helpers/entities/service/form';
import { patternToForm } from '@/helpers/entities/pattern/form';
import { infosToArray } from '@/helpers/entities/shared/form';

import ServiceForm from '@/components/other/service/form/service-form.vue';

const stubs = {
  'c-name-field': true,
  'c-entity-category-field': true,
  'c-alarm-state-field': true,
  'c-impact-level-field': true,
  'c-coordinates-field': true,
  'text-editor-field': true,
  'c-enabled-field': true,
  'entity-state-setting': true,
  'c-patterns-field': true,
  'manage-infos': true,
};

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectCategoryField = wrapper => wrapper.find('c-entity-category-field-stub');
const selectAlarmStateField = wrapper => wrapper.find('c-alarm-state-field-stub');
const selectImpactLevelField = wrapper => wrapper.find('c-impact-level-field-stub');
const selectCoordinatesField = wrapper => wrapper.find('c-coordinates-field-stub');
const selectTextEditorField = wrapper => wrapper.find('text-editor-field-stub');
const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectPatternsField = wrapper => wrapper.find('c-patterns-field-stub');
const selectManageInfosField = wrapper => wrapper.find('manage-infos-stub');
const selectEntityStateSettingField = wrapper => wrapper.find('entity-state-setting-stub');

describe('service-form', () => {
  const factory = generateShallowRenderer(ServiceForm, { stubs });
  const snapshotFactory = generateRenderer(ServiceForm, { stubs });

  const defaultServiceForm = serviceToForm();

  test('Name changed after trigger name field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newName = Faker.datatype.string();

    selectNameField(wrapper).triggerCustomEvent('input', newName);

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      name: newName,
    });
  });

  test('Category changed after trigger category field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newCategory = {
      _id: Faker.datatype.string(),
    };

    selectCategoryField(wrapper).triggerCustomEvent('input', newCategory);

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      category: newCategory,
    });
  });

  test('Available state changed after trigger alarm state field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newState = ALARM_STATES.minor;

    selectAlarmStateField(wrapper).triggerCustomEvent('input', newState);

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      sli_avail_state: newState,
    });
  });

  test('Impact level changed after trigger impact level field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newImpactLevel = ALARM_STATES.minor;

    selectImpactLevelField(wrapper).triggerCustomEvent('input', newImpactLevel);

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      impact_level: newImpactLevel,
    });
  });

  test('Coordinates changed after trigger coordinates field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newCoordinates = {
      lat: Faker.datatype.number(),
      lng: Faker.datatype.number(),
    };

    selectCoordinatesField(wrapper).triggerCustomEvent('input', newCoordinates);

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      coordinates: newCoordinates,
    });
  });

  test('Output template changed after trigger text editor field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newTemplate = Faker.datatype.string();

    selectTextEditorField(wrapper).triggerCustomEvent('input', newTemplate);

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      output_template: newTemplate,
    });
  });

  test('Enabled changed after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newEnabled = Faker.datatype.boolean();

    selectEnabledField(wrapper).triggerCustomEvent('input', newEnabled);

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      enabled: newEnabled,
    });
  });

  test('Prepare function passed to state setting', () => {
    const prepareStateSettingForm = jest.fn();
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
        prepareStateSettingForm,
      },
    });

    expect(
      selectEntityStateSettingField(wrapper).vm.preparer,
    ).toBe(prepareStateSettingForm);
  });

  test('Patterns changed after trigger patterns field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newPatterns = patternToForm({
      entity_pattern: [[
        {
          field: ENTITY_PATTERN_FIELDS.name,
          cond: {
            type: PATTERN_CONDITIONS.notEqual,
            value: 'test',
          },
        },
      ]],
    });

    selectPatternsField(wrapper).triggerCustomEvent('input', newPatterns);

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      patterns: newPatterns,
    });
  });

  test('Infos changed after trigger manage infos field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newInfos = infosToArray({
      [Faker.datatype.string()]: {
        value: Faker.datatype.string(),
        description: Faker.datatype.string(),
      },
    });

    selectManageInfosField(wrapper).triggerCustomEvent('input', newInfos);

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      infos: newInfos,
    });
  });

  test('Renders `service-form` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `service-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          ...defaultServiceForm,
          impact_level: 2,
          name: 'service-form',
          category: {
            _id: 'category-id',
          },
          enabled: false,
          output_template: 'output-template',
          sli_avail_state: ALARM_STATES.critical,
          coordinates: {
            lat: 2,
            lng: 3,
          },
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
