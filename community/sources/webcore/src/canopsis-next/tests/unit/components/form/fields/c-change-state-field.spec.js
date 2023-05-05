import Faker from 'faker';
import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { createMockedStoreGetters } from '@unit/utils/store';
import { createTextareaInputStub, createNumberInputStub } from '@unit/stubs/input';
import flushPromises from 'flush-promises';
import { ENTITIES_STATES } from '@/constants';
import CChangeStateField from '@/components/forms/fields/c-change-state-field.vue';

const stubs = {
  'state-criticity-field': createNumberInputStub('state-criticity-field'),
  'v-textarea': createTextareaInputStub('v-textarea'),
};

const snapshotStubs = {
  'state-criticity-field': true,
};

const factory = generateShallowRenderer(CChangeStateField, { stubs,
});

const snapshotFactory = generateRenderer(CChangeStateField, {

  stubs: snapshotStubs,
  parentComponent: {
    $_veeValidate: {
      validator: 'new',
    },
  },
});

describe('c-change-state-field', () => {
  it('State changed after trigger the state field', () => {
    const initialValue = {
      state: ENTITIES_STATES.major,
      output: Faker.datatype.string(),
    };
    const wrapper = factory({
      store: createMockedStoreGetters({ name: 'info', allowChangeSeverityToInfo: false }),
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
      store: createMockedStoreGetters({ name: 'info', allowChangeSeverityToInfo: false }),
      propsData: {
        value: initialValue,
      },
    });

    const outputElement = wrapper.find('.v-textarea textarea');
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

  it('Renders `c-change-state-field` after validate correctly', async () => {
    const value = {
      state: ENTITIES_STATES.major,
      output: '',
    };

    const wrapper = snapshotFactory({
      store: createMockedStoreGetters({ name: 'info', allowChangeSeverityToInfo: false }),
      propsData: {
        value,
      },
    });

    const { $validator: validator } = wrapper.vm;

    await validator.validateAll();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-change-state-field` with custom label correctly', () => {
    const wrapper = snapshotFactory({
      store: createMockedStoreGetters({ name: 'info', allowChangeSeverityToInfo: false }),
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
    const wrapper = snapshotFactory({
      store: createMockedStoreGetters({ name: 'info', allowChangeSeverityToInfo: false }),
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
    const wrapper = snapshotFactory({
      store: createMockedStoreGetters({ name: 'info', allowChangeSeverityToInfo: true }),
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

  it('Renders `c-change-state-field` with errors correctly', async () => {
    const name = 'customName';

    const wrapper = snapshotFactory({
      store: createMockedStoreGetters({ name: 'info', allowChangeSeverityToInfo: true }),
      propsData: {
        value: {
          state: ENTITIES_STATES.ok,
          output: 'Output',
        },
        name,
      },
    });

    const { $validator: validator } = wrapper.vm;
    validator.errors.add([
      {
        field: `${name}.output`,
        msg: 'Output error',
      },
    ]);

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
