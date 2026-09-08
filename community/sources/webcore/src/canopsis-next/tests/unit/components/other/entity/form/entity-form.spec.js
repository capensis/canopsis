import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { getFormGeneralPatternsTabsStub } from '@unit/stubs/form';

import { ALARM_STATES, BASIC_ENTITY_TYPES, ENTITY_TYPES } from '@/constants';

import { infosToArray } from '@/helpers/entities/shared/form';
import { entityToForm } from '@/helpers/entities/entity/form';

import EntityForm from '@/components/other/entity/form/entity-form.vue';

const stubs = {
  'entity-general-form': true,
  'entity-manage-infos-form': true,
  'c-form-general-patterns-tabs': getFormGeneralPatternsTabsStub(),
};

const selectEntityGeneralForm = wrapper => wrapper.find('entity-general-form-stub');
const selectEntityManageInfosForm = wrapper => wrapper.find('entity-manage-infos-form-stub');

describe('entity-form', () => {
  const factory = generateShallowRenderer(EntityForm, { stubs });
  const snapshotFactory = generateRenderer(EntityForm, { stubs });

  const defaultEntityForm = entityToForm();

  test('General form is rendered in general tab', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    expect(selectEntityGeneralForm(wrapper).exists()).toBe(true);
  });

  test('Manage infos form is rendered in additional tab', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    expect(selectEntityManageInfosForm(wrapper).exists()).toBe(true);
  });

  test('Name changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newName = Faker.datatype.string();

    selectEntityGeneralForm(wrapper).triggerCustomEvent('input', {
      ...defaultEntityForm,
      name: newName,
    });

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      name: newName,
    });
  });

  test('Description changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newDescription = Faker.datatype.string();

    selectEntityGeneralForm(wrapper).triggerCustomEvent('input', {
      ...defaultEntityForm,
      description: newDescription,
    });

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      description: newDescription,
    });
  });

  test('Available state changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newState = ALARM_STATES.minor;

    selectEntityGeneralForm(wrapper).triggerCustomEvent('input', {
      ...defaultEntityForm,
      sli_avail_state: newState,
    });

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      sli_avail_state: newState,
    });
  });

  test('Entity type changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newType = BASIC_ENTITY_TYPES.connector;

    selectEntityGeneralForm(wrapper).triggerCustomEvent('input', {
      ...defaultEntityForm,
      type: newType,
    });

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      type: newType,
    });
  });

  test('Impact level changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newImpactLevel = ALARM_STATES.minor;

    selectEntityGeneralForm(wrapper).triggerCustomEvent('input', {
      ...defaultEntityForm,
      impact_level: newImpactLevel,
    });

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      impact_level: newImpactLevel,
    });
  });

  test('Coordinates changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newCoordinates = {
      lat: Faker.datatype.number(),
      lng: Faker.datatype.number(),
    };

    selectEntityGeneralForm(wrapper).triggerCustomEvent('input', {
      ...defaultEntityForm,
      coordinates: newCoordinates,
    });

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      coordinates: newCoordinates,
    });
  });

  test('Enabled changed after trigger general form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newEnabled = Faker.datatype.boolean();

    selectEntityGeneralForm(wrapper).triggerCustomEvent('input', {
      ...defaultEntityForm,
      enabled: newEnabled,
    });

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      enabled: newEnabled,
    });
  });

  test('Prepare function passed to general form', () => {
    const prepareStateSettingForm = jest.fn();
    const wrapper = factory({
      propsData: {
        form: {
          ...defaultEntityForm,
          type: ENTITY_TYPES.component,
        },
        prepareStateSettingForm,
      },
    });

    expect(selectEntityGeneralForm(wrapper).props('prepareStateSettingForm')).toBe(prepareStateSettingForm);
  });

  test('Infos changed after trigger manage infos form', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newInfos = infosToArray({
      [Faker.datatype.string()]: {
        value: Faker.datatype.string(),
        description: Faker.datatype.string(),
      },
    });

    selectEntityManageInfosForm(wrapper).triggerCustomEvent('input', newInfos);

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      infos: newInfos,
    });
  });

  test('Renders `entity-form` with default form', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `entity-form` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          ...defaultEntityForm,
          name: 'entity-form',
          description: 'entity-form-description',
          type: BASIC_ENTITY_TYPES.resource,
          enabled: false,
          impact_level: 2,
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
