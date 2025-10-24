import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';
import { createSelectInputStub } from '@unit/stubs/input';

import { ALARM_PATTERN_FIELDS } from '@/constants';

import PatternAttributeField from '@/components/forms/fields/pattern/pattern-attribute-field.vue';
import CSelectField from '@/components/forms/fields/c-select-field';

const stubs = {
  'c-select-field': createSelectInputStub('c-select-field'),
  'c-simple-tooltip': true,
};

const snapshotStubs = {
  'c-select-field': CSelectField,
  'c-simple-tooltip': true,
};

const selectSelectField = wrapper => wrapper.find('.c-select-field');

describe('pattern-attribute-field', () => {
  const factory = generateShallowRenderer(PatternAttributeField, { stubs });
  const snapshotFactory = generateRenderer(PatternAttributeField, { stubs: snapshotStubs });

  it('Value changed after trigger the input', () => {
    const value = {
      value: ALARM_PATTERN_FIELDS.ack,
      text: 'Text',
    };
    const wrapper = factory({
      propsData: {
        value,
      },
    });
    const selectElement = selectSelectField(wrapper);

    selectElement.triggerCustomEvent('input', value);

    expect(wrapper).toEmitInput(value);
  });

  it('Renders `pattern-attribute-field` with default props', async () => {
    const value = {
      value: ALARM_PATTERN_FIELDS.component,
      text: 'Component',
    };
    const wrapper = snapshotFactory({
      propsData: {
        value,
      },
    });

    await wrapper.activateAllMenus();

    expect(wrapper).toMatchSnapshot();
    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `pattern-attribute-field` with custom props', async () => {
    const value = {
      value: ALARM_PATTERN_FIELDS.component,
      text: 'Component',
    };
    const wrapper = snapshotFactory({
      propsData: {
        value,
        label: 'Custom label',
        items: [
          {
            value: ALARM_PATTERN_FIELDS.ackAt,
            text: 'Ack at',
          },
        ],
        name: 'custom_filter_attribute_name',
        disabled: true,
      },
    });

    expect(wrapper).toMatchSnapshot();

    await wrapper.activateAllMenus();

    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `pattern-attribute-field` menu with alias item selected', async () => {
    const items = [
      {
        value: ALARM_PATTERN_FIELDS.component,
        originalValue: 'component',
        text: 'Component',
      },
      {
        value: ALARM_PATTERN_FIELDS.infos,
        originalValue: 'infos',
        text: 'Infos',
        options: { alias: true },
      },
    ];

    const wrapper = snapshotFactory({
      propsData: {
        value: ALARM_PATTERN_FIELDS.infos,
        items,
        name: 'alias_attribute',
      },
    });

    expect(wrapper).toMatchSnapshot();

    await wrapper.activateAllMenus();

    expect(wrapper).toMatchMenuSnapshot();
  });

  it('Renders `pattern-attribute-field` menu with alias item in list', async () => {
    const items = [
      {
        value: ALARM_PATTERN_FIELDS.component,
        originalValue: 'component',
        text: 'Component',
      },
      {
        value: ALARM_PATTERN_FIELDS.infos,
        originalValue: 'infos',
        text: 'Infos',
        options: { alias: true },
      },
    ];

    const wrapper = snapshotFactory({
      propsData: {
        value: ALARM_PATTERN_FIELDS.component,
        items,
        name: 'alias_attribute_menu',
      },
    });

    expect(wrapper).toMatchSnapshot();

    await wrapper.activateAllMenus();

    expect(wrapper).toMatchMenuSnapshot();
  });
});
