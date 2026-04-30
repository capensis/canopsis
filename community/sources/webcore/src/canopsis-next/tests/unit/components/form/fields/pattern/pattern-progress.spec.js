import { flushPromises, generateRenderer } from '@unit/utils/vue';

import PatternProgress from '@/components/forms/fields/pattern/pattern-progress.vue';

const progressOverlayStub = {
  name: 'CProgressOverlay',
  props: ['pending', 'opacity', 'size', 'width'],
  template: `
    <div class="c-progress-overlay-stub">
      <slot name="progress" />
      <slot />
    </div>
  `,
};

const stubs = {
  'c-progress-overlay': progressOverlayStub,
};

describe('pattern-progress', () => {
  const factory = generateRenderer(PatternProgress, { stubs });

  test('Should emit close when close button is clicked in failed state', async () => {
    const wrapper = factory({
      propsData: {
        failedReason: 'network error',
      },
    });

    await flushPromises();

    const buttons = wrapper.findAllComponents({ name: 'VBtn' });

    await buttons.at(0).trigger('click');
    await flushPromises();

    expect(wrapper).toHaveBeenEmit('close');
  });

  test('Should emit try-again when try again button is clicked in failed state', async () => {
    const wrapper = factory({
      propsData: {
        failedReason: 'network error',
      },
    });

    await flushPromises();

    const buttons = wrapper.findAllComponents({ name: 'VBtn' });

    await buttons.at(1).trigger('click');
    await flushPromises();

    expect(wrapper).toHaveBeenEmit('try-again');
  });

  test('Should emit cancel when cancel button is clicked in progress state', async () => {
    const wrapper = factory({
      propsData: {
        failedReason: '',
      },
    });

    await flushPromises();

    const button = wrapper.findComponent({ name: 'VBtn' });

    await button.trigger('click');
    await flushPromises();

    expect(wrapper).toHaveBeenEmit('cancel');
  });

  test('Should render custom labels and message when corresponding props are set', async () => {
    const wrapper = factory({
      propsData: {
        failedReason: 'x',
        failedMessageText: 'Custom failure',
        inProgressText: 'Custom progress',
        closeButtonLabel: 'Dismiss',
        tryAgainButtonLabel: 'Retry',
        cancelButtonLabel: 'Stop',
      },
    });

    await flushPromises();

    expect(wrapper.text()).toContain('Custom failure');

    wrapper.setProps({ failedReason: '' });
    await flushPromises();

    expect(wrapper.text()).toContain('Custom progress');
    expect(wrapper.text()).toContain('Stop');
  });

  test('Renders `pattern-progress` in progress state', () => {
    const wrapper = factory({
      propsData: {
        failedReason: '',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pattern-progress` in failed state', () => {
    const wrapper = factory({
      propsData: {
        failedReason: 'test reason',
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
