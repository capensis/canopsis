import { generateShallowRenderer, generateRenderer } from '@unit/utils/vue';

import { ALARM_STATES } from '@/constants';

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
        value: ALARM_STATES.major,
      },
    });

    wrapper.find('input.v-btn-toggle').setValue(ALARM_STATES.ok);

    expect(wrapper).toEmitInput(ALARM_STATES.ok);
  });

  it('Renders `state-criticity-field` with default props correctly', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `state-criticity-field` with custom props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ALARM_STATES.major,
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
