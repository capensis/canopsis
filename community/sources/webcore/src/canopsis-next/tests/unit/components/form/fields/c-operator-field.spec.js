import { mount, createVueInstance } from '@unit/utils/vue';

import COperatorField from '@/components/forms/fields/c-operator-field.vue';
import { FILTER_MONGO_OPERATORS } from '@/constants';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(COperatorField, {
  localVue,

  ...options,
});

const selectRadioElementByValue = (wrapper, value) => wrapper.find(`.v-input--radio-group__input input[value="${value}"]`);

describe('c-operator-field', () => {
  it.each([FILTER_MONGO_OPERATORS.and, FILTER_MONGO_OPERATORS.or])(
    'Value changed to %s after trigger radio button',
    (value) => {
      const wrapper = snapshotFactory({
        propsData: {
          value: value === FILTER_MONGO_OPERATORS.and
            ? FILTER_MONGO_OPERATORS.or
            : FILTER_MONGO_OPERATORS.and,
        },
      });
      const radioElement = selectRadioElementByValue(wrapper, value);

      radioElement.trigger('change');

      const inputEvents = wrapper.emitted('input');

      expect(inputEvents).toHaveLength(1);

      const [eventData] = inputEvents[0];
      expect(eventData).toEqual(value);
    },
  );

  it('Renders `c-operator-field` with default props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: FILTER_MONGO_OPERATORS.or,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `c-operator-field` with custom props correctly', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: FILTER_MONGO_OPERATORS.and,
        name: 'customName',
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
