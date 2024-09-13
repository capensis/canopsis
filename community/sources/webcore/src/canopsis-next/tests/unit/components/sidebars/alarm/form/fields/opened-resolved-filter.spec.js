import { generateRenderer } from '@unit/utils/vue';

import { ALARMS_OPENED_VALUES } from '@/constants';

import OpenedResolvedFilter from '@/components/sidebars/alarm/form/fields/opened-resolved-filter.vue';

const selectRadioElementsByValue = (wrapper, value) => wrapper
  .find(`.v-input--radio-group__input input[value="${String(value)}"]`);

describe('opened-resolved-filter', () => {
  const snapshotFactory = generateRenderer(OpenedResolvedFilter, {
    parentComponent: {
      provide: {
        list: {
          register: jest.fn(),
          unregister: jest.fn(),
        },
        listClick: jest.fn(),
      },
    },
  });

  it.each(Object.entries(ALARMS_OPENED_VALUES))(
    'Value changed after to %s after trigger a radio field',
    (key, value) => {
      const wrapper = snapshotFactory({
        propsData: {
          value: value === ALARMS_OPENED_VALUES.opened
            ? ALARMS_OPENED_VALUES.resolved
            : ALARMS_OPENED_VALUES.opened,
        },
      });

      selectRadioElementsByValue(wrapper, value).trigger('change');

      expect(wrapper).toEmitInput(value);
    },
  );

  it('Renders `opened-resolved-filter` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  it('Renders `opened-resolved-filter` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: ALARMS_OPENED_VALUES.all,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
