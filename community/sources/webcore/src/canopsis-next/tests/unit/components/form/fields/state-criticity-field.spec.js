import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { ENTITIES_STATES } from '@/constants';
import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

const stubs = {
  'v-btn-toggle': {
    props: ['value'],
    template: `
     <input
       :value="value"
       class="v-btn-toggle"
       @input="$listeners.change(+$event.target.value)"
     />
   `,
  },
};

describe('state-criticity-field', () => {
  const factory = generateShallowRenderer(StateCriticityField, { stubs });
  const snapshotFactory = generateRenderer(StateCriticityField);

  it('Value changed after trigger click on the button', () => {
    const wrapper = factory({
      propsData: {
        value: ENTITIES_STATES.major,
      },
    });

    const buttonToggleElement = wrapper.find('input.v-btn-toggle');

    buttonToggleElement.setValue(ENTITIES_STATES.ok);

    const inputEvents = wrapper.emitted('input');
    expect(inputEvents).toHaveLength(1);

    const [eventData] = inputEvents[0];
    expect(eventData).toEqual(ENTITIES_STATES.ok);
  });

  it('Renders `state-criticity-field` with default props correctly', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `state-criticity-field` with custom props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ENTITIES_STATES.major,
        mandatory: true,
        stateValues: {
          ok: 0,
          minor: 1,
          critical: 3,
        },
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
