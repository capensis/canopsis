import { getFormGeneralPatternsTabsStub } from '@unit/stubs/form';
import { createMockedStoreModules, createPatternsFieldsModule } from '@unit/utils/store';
import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import { ALARM_STATES, ENTITY_PATTERN_FIELDS, PATTERN_CONDITIONS } from '@/constants';

import { patternToForm } from '@/helpers/entities/pattern/form';
import { serviceToForm } from '@/helpers/entities/service/form';
import { infosToArray } from '@/helpers/entities/shared/form';

import ServiceForm from '@/components/other/service/form/service-form.vue';

const serviceGeneralFormStub = {
  name: 'ServiceGeneralForm',
  props: ['form', 'prepareStateSettingForm', 'templateVars'],
  emits: ['update:form'],
  template: '<div />',
};

const stubs = {
  'c-enabled-field': true,
  'service-general-form': serviceGeneralFormStub,
  'service-entity-patterns-form': true,
  'service-manage-infos-form': true,
  'c-form-general-patterns-tabs': getFormGeneralPatternsTabsStub(),
};

const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectServiceGeneralForm = wrapper => wrapper.findComponent(serviceGeneralFormStub);
const selectServiceEntityPatternsForm = wrapper =>
  wrapper.find('service-entity-patterns-form-stub');
const selectServiceManageInfosForm = wrapper => wrapper.find('service-manage-infos-form-stub');

describe('service-form', () => {
  const { patternsFieldsModule, fetchServicePatternFields } = createPatternsFieldsModule();

  fetchServicePatternFields.mockResolvedValue({
    entity_pattern: [
      { name: 'name', enabled: true, alias: false },
      { name: 'category', enabled: true, alias: false },
      { name: 'component', enabled: true, alias: false },
    ],
    alarm_pattern: [],
    event_pattern: [],
  });

  const store = createMockedStoreModules([patternsFieldsModule]);

  const factory = generateShallowRenderer(ServiceForm, { stubs, store });

  const defaultServiceForm = serviceToForm();

  test('General form is rendered in general tab', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    expect(selectServiceGeneralForm(wrapper).exists()).toBe(true);
  });

  test('Entity patterns form is rendered in additional tab', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    expect(selectServiceEntityPatternsForm(wrapper).exists()).toBe(true);
  });

  test('Manage infos form is rendered in patterns tab', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    expect(selectServiceManageInfosForm(wrapper).exists()).toBe(true);
  });

  test('Name changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newName = Faker.datatype.string();

    selectServiceGeneralForm(wrapper).triggerCustomEvent('update:form', {
      ...defaultServiceForm,
      name: newName,
    });

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      name: newName,
    });
  });

  test('Category changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newCategory = {
      _id: Faker.datatype.string(),
    };

    selectServiceGeneralForm(wrapper).triggerCustomEvent('update:form', {
      ...defaultServiceForm,
      category: newCategory,
    });

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      category: newCategory,
    });
  });

  test('Available state changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newState = ALARM_STATES.minor;

    selectServiceGeneralForm(wrapper).triggerCustomEvent('update:form', {
      ...defaultServiceForm,
      sli_avail_state: newState,
    });

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      sli_avail_state: newState,
    });
  });

  test('Impact level changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newImpactLevel = ALARM_STATES.minor;

    selectServiceGeneralForm(wrapper).triggerCustomEvent('update:form', {
      ...defaultServiceForm,
      impact_level: newImpactLevel,
    });

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      impact_level: newImpactLevel,
    });
  });

  test('Coordinates changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newCoordinates = {
      lat: Faker.datatype.number(),
      lng: Faker.datatype.number(),
    };

    selectServiceGeneralForm(wrapper).triggerCustomEvent('update:form', {
      ...defaultServiceForm,
      coordinates: newCoordinates,
    });

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      coordinates: newCoordinates,
    });
  });

  test('Output template changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newTemplate = Faker.datatype.string();

    selectServiceGeneralForm(wrapper).triggerCustomEvent('update:form', {
      ...defaultServiceForm,
      output_template: newTemplate,
    });

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

  test('Prepare function passed to general form', () => {
    const prepareStateSettingForm = vi.fn();
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
        prepareStateSettingForm,
      },
    });

    expect(selectServiceGeneralForm(wrapper).props('prepareStateSettingForm')).toBe(
      prepareStateSettingForm,
    );
  });

  test('Patterns changed after trigger entity patterns form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultServiceForm,
      },
    });

    const newPatterns = patternToForm({
      entity_pattern: [
        [
          {
            field: ENTITY_PATTERN_FIELDS.name,
            cond: {
              type: PATTERN_CONDITIONS.notEqual,
              value: 'test',
            },
          },
        ],
      ],
    });

    selectServiceEntityPatternsForm(wrapper).triggerCustomEvent('input', newPatterns);

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      patterns: newPatterns,
    });
  });

  test('Infos changed after trigger manage infos form', () => {
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

    selectServiceManageInfosForm(wrapper).triggerCustomEvent('update:infos', newInfos);

    expect(wrapper).toEmitInput({
      ...defaultServiceForm,
      infos: newInfos,
    });
  });
});
