import { generateRenderer } from '@unit/utils/vue';

import ExtraDetailsCloseDelay
  from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-close-delay.vue';

const stubs = {
  'c-alarm-extra-details-chip': true,
  'c-simple-tooltip': true,
};

describe('extra-details-close-delay', () => {
  const closeDelay = {
    _t: 'statedec',
    t: 1743511657,
    a: 'system',
    user_id: '',
    m: 'closed after 20 seconds delay',
    val: 0,
    initiator: 'system',
  };
  const snapshotFactory = generateRenderer(ExtraDetailsCloseDelay, {
    stubs,
    attachTo: document.body,
  });

  it('Renders `extra-details-close-delay` with full children and rule', () => {
    jest.useFakeTimers();

    const wrapper = snapshotFactory({
      propsData: {
        closeDelay,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
