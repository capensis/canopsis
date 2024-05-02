import Faker from 'faker';

import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { ALARM_STATES, BASIC_ENTITY_TYPES, ENTITY_TYPES } from '@/constants';

import { infosToArray } from '@/helpers/entities/shared/form';
import { entityToForm } from '@/helpers/entities/entity/form';

import EntityForm from '@/components/other/entity/form/entity-form.vue';

const stubs = {
  'c-name-field': true,
  'c-description-field': true,
  'c-enabled-field': true,
  'c-impact-level-field': true,
  'c-alarm-state-field': true,
  'c-entity-type-field': true,
  'c-coordinates-field': true,
  'entity-state-setting': true,
  'manage-infos': true,
};

const selectNameField = wrapper => wrapper.find('c-name-field-stub');
const selectDescriptionField = wrapper => wrapper.find('c-description-field-stub');
const selectEnabledField = wrapper => wrapper.find('c-enabled-field-stub');
const selectImpactLevelField = wrapper => wrapper.find('c-impact-level-field-stub');
const selectAlarmStateField = wrapper => wrapper.find('c-alarm-state-field-stub');
const selectEntityTypeField = wrapper => wrapper.find('c-entity-type-field-stub');
const selectCoordinatesField = wrapper => wrapper.find('c-coordinates-field-stub');
const selectEntityStateSettingField = wrapper => wrapper.find('entity-state-setting-stub');
const selectManageInfosField = wrapper => wrapper.find('manage-infos-stub');

describe('entity-form', () => {
  const factory = generateShallowRenderer(EntityForm, { stubs });
  const snapshotFactory = generateRenderer(EntityForm, { stubs });

  const defaultEntityForm = entityToForm();

  test('Name changed after trigger name field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newName = Faker.datatype.string();

    selectNameField(wrapper).triggerCustomEvent('input', newName);

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      name: newName,
    });
  });

  test('Description changed after trigger description field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newDescription = Faker.datatype.string();

    selectDescriptionField(wrapper).triggerCustomEvent('input', newDescription);

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      description: newDescription,
    });
  });

  test('Available state changed after trigger alarm state field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newState = ALARM_STATES.minor;

    selectAlarmStateField(wrapper).triggerCustomEvent('input', newState);

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      sli_avail_state: newState,
    });
  });

  test('Entity type changed after trigger entity type field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newType = BASIC_ENTITY_TYPES.connector;

    selectEntityTypeField(wrapper).triggerCustomEvent('input', newType);

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      type: newType,
    });
  });

  test('Impact level changed after trigger impact level field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newImpactLevel = ALARM_STATES.minor;

    selectImpactLevelField(wrapper).triggerCustomEvent('input', newImpactLevel);

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      impact_level: newImpactLevel,
    });
  });

  test('Coordinates changed after trigger coordinates field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newCoordinates = {
      lat: Faker.datatype.number(),
      lng: Faker.datatype.number(),
    };

    selectCoordinatesField(wrapper).triggerCustomEvent('input', newCoordinates);

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      coordinates: newCoordinates,
    });
  });

  test('Enabled changed after trigger enabled field', () => {
    const wrapper = factory({
      propsData: {
        form: defaultEntityForm,
      },
    });

    const newEnabled = Faker.datatype.boolean();

    selectEnabledField(wrapper).triggerCustomEvent('input', newEnabled);

    expect(wrapper).toEmitInput({
      ...defaultEntityForm,
      enabled: newEnabled,
    });
  });

  test('Prepare function passed to state setting', () => {
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

    expect(
      selectEntityStateSettingField(wrapper).vm.preparer,
    ).toBe(prepareStateSettingForm);
  });

  test('Infos changed after trigger manage infos field', () => {
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

    selectManageInfosField(wrapper).triggerCustomEvent('input', newInfos);

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
