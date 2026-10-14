import { mockModals, mockSocket } from '@unit/utils/mock-hooks';
import { generateShallowRenderer } from '@unit/utils/vue';
import { flushPromises } from '@vue/test-utils';

import { MODALS } from '@/constants';

import PbehaviorsListAlarmFilteringBtn from '@/components/other/pbehavior/pbehaviors/partials/pbehaviors-list-alarm-filtering-btn.vue';

const { runAlarmFiltering } = vi.hoisted(() => ({
  runAlarmFiltering: vi.fn().mockResolvedValue(),
}));

vi.mock('@/hooks/store/modules/pbehavior-patterns', () => ({
  usePbehaviorPatterns: () => ({ runAlarmFiltering }),
}));

const stubs = {
  'v-progress-circular': true,
};

const selectButton = wrapper => wrapper.find('v-btn-stub');

describe('pbehaviors-list-alarm-filtering-btn', () => {
  const $modals = mockModals();
  const $socket = mockSocket();

  const factory = generateShallowRenderer(PbehaviorsListAlarmFilteringBtn, {
    stubs,
    mocks: {
      $modals,
      $socket,
    },
  });

  afterEach(() => {
    $modals.show.mockClear();
    runAlarmFiltering.mockClear();
  });

  test('Should render button with correct structure', async () => {
    const wrapper = factory();

    await flushPromises();

    const button = selectButton(wrapper);
    expect(button.exists()).toBe(true);
    expect(button.attributes('color')).toBe('primary');
  });

  test('Should have required component structure', async () => {
    const wrapper = factory();

    await flushPromises();

    // Check basic component structure
    expect(wrapper.exists()).toBe(true);
    expect(selectButton(wrapper).exists()).toBe(true);

    // Check that component has the expected elements
    const button = selectButton(wrapper);
    expect(button.attributes('color')).toBe('primary');
  });

  test('Should handle component lifecycle correctly', async () => {
    const wrapper = factory();

    await flushPromises();

    // Component should mount without errors
    expect(wrapper.vm).toBeDefined();
    expect(wrapper.element).toBeDefined();
  });

  test('Should trigger modal when button is clicked', async () => {
    const wrapper = factory();

    await selectButton(wrapper).trigger('click');

    expect($modals.show).toHaveBeenCalledWith({
      name: MODALS.confirmation,
      config: {
        title: expect.any(String),
        text: expect.any(String),
        action: expect.any(Function),
      },
    });
  });

  test('Should call runAlarmFiltering when modal action is executed', async () => {
    const wrapper = factory();

    await selectButton(wrapper).trigger('click');

    // Get the modal config from the first call
    const [{ config }] = $modals.show.mock.calls[0];

    // Execute the action from the modal config
    await config.action();

    expect(runAlarmFiltering).toHaveBeenCalled();
  });
});
