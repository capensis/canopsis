import Faker from 'faker';
import { Validator } from 'vee-validate';
import { mount, shallowMount, createVueInstance } from '@unit/utils/vue';

import { createMockedStoreGetters } from '@unit/utils/store';
import { createInputStub, createNumberInputStub } from '@unit/stubs/input';
import { ENTITIES_STATES } from '@/constants';
import CChangeStateField from '@/components/forms/fields/c-change-state-field.vue';

const localVue = createVueInstance();

const stubs = {
  'state-criticity-field': createNumberInputStub('state-criticity-field'),
  'v-textarea': createInputStub('v-textarea'),
};

const snapshotStubs = {
};

const factory = (options = {}) => shallowMount(CChangeStateField, {
  localVue,
  stubs,
  ...options,
});

describe('c-change-state-field', () => {
  it('State changed after trigger the state field', () => {
    const initialValue = {
      state: ENTITIES_STATES.major,
      output: Faker.datatype.string(),
    };
    const wrapper = factory({
      store: createMockedStoreGetters('info', { allowChangeSeverityToInfo: false }),
      propsData: {
        value: initialValue,
      },
    });

    const stateCriticityElement = wrapper.find('input.state-criticity-field');

    stateCriticityElement.setValue(ENTITIES_STATES.critical);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData.state).toEqual(ENTITIES_STATES.critical);
    expect(eventData).toEqual({
      ...initialValue,
      state: ENTITIES_STATES.critical,
    });
  });

  it('Output changed after trigger the textarea field', () => {
    const initialValue = {
      state: ENTITIES_STATES.major,
      output: Faker.datatype.string(),
    };
    const wrapper = factory({
      store: createMockedStoreGetters('info', { allowChangeSeverityToInfo: false }),
      propsData: {
        value: initialValue,
      },
    });

    const outputElement = wrapper.find('input.v-textarea');

    const newOutput = Faker.datatype.string();
    outputElement.setValue(newOutput);

    const inputEvents = wrapper.emitted('input');

    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData.output).toEqual(newOutput);
    expect(eventData).toEqual({
      ...initialValue,
      output: newOutput,
    });
  });

  it('Error created for output field, after validate', async () => {
    const validator = new Validator();
    const name = Faker.datatype.string();
    const value = {
      state: ENTITIES_STATES.major,
      output: '',
    };

    factory({
      provide: {
        $validator: validator,
      },
      store: createMockedStoreGetters('info', { allowChangeSeverityToInfo: false }),
      propsData: {
        value,
        name,
      },
    });

    const isValid = await validator.validateAll();

    expect(isValid).toBe(false);
    expect(validator.errors.count()).toBe(1);

    const [error] = validator.errors;
    expect(error.field).toBe(`${name}.output`);
    expect(error.rule).toBe('required');
  });

  it('Renders `c-change-state-field` with custom label correctly', () => {
    const wrapper = mount(CChangeStateField, {
      localVue,
      stubs: snapshotStubs,
      store: createMockedStoreGetters('info', { allowChangeSeverityToInfo: false }),
      propsData: {
        value: {
          state: ENTITIES_STATES.ok,
          output: 'Custom label output',
        },
        name: 'customLabelFieldName',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-change-state-field` without allowed change severity to info correctly', () => {
    const wrapper = mount(CChangeStateField, {
      localVue,
      stubs: snapshotStubs,
      store: createMockedStoreGetters('info', { allowChangeSeverityToInfo: false }),
      propsData: {
        value: {
          state: ENTITIES_STATES.ok,
          output: 'Output',
        },
        label: 'Custom label',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-change-state-field` with allowed change severity to info correctly', () => {
    const wrapper = mount(CChangeStateField, {
      localVue,
      stubs: snapshotStubs,
      store: createMockedStoreGetters('info', { allowChangeSeverityToInfo: true }),
      propsData: {
        value: {
          state: ENTITIES_STATES.ok,
          output: 'Output',
        },
        label: 'Custom label',
        name: 'customLabelFieldName',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-change-state-field` with errors correctly', () => {
    const name = 'customName';
    const validator = new Validator();
    validator.errors.add([
      {
        field: `${name}.output`,
        msg: 'Output error',
      },
    ]);

    const wrapper = mount(CChangeStateField, {
      localVue,
      stubs: snapshotStubs,
      provide: {
        $validator: validator,
      },
      store: createMockedStoreGetters('info', { allowChangeSeverityToInfo: true }),
      propsData: {
        value: {
          state: ENTITIES_STATES.ok,
          output: 'Output',
        },
        name,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
