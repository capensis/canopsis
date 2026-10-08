import { generateShallowRenderer } from '@unit/utils/vue';

import ExtraDetailsPause from '@/components/widgets/alarm/columns-formatting/extra-details/extra-details-pause.vue';

const CTooltipStub = {
  name: 'c-tooltip',
  props: ['content'],
  template: '<div><slot name="activator" :props="{}" /></div>',
};

describe('extra-details-pause', () => {
  const factory = generateShallowRenderer(ExtraDetailsPause, {
    stubs: {
      'c-tooltip': CTooltipStub,
      'c-alarm-extra-details-chip': true,
    },
  });

  it('shows only the pause comment in the tooltip', () => {
    const wrapper = factory({
      propsData: {
        pause: {
          m: 'Maintenance window',
          pause_type: 'maintenance',
          author: 'John Doe',
          timestamp: 1641768553,
        },
      },
    });
    const content = wrapper.findComponent(CTooltipStub).props('content');

    expect(content).toContain('Maintenance window');
    expect(content).not.toContain('maintenance');
    expect(content).not.toContain('<strong>Pause</strong>');
    expect(content).not.toContain('Comment:');
    expect(content).not.toContain('Type:');
    expect(content).not.toContain('John Doe');
    expect(content).not.toContain('1641768553');
  });

  it('does not show a tooltip without a comment', () => {
    const wrapper = factory({ propsData: { pause: {} } });

    expect(wrapper.findComponent(CTooltipStub).exists()).toBe(false);
  });

  it('shows the inherited message when entity and alarm pause comments match', () => {
    const wrapper = factory({
      propsData: {
        pause: {
          m: 'Maintenance window',
          pause_type: 'maintenance',
        },
        entity: {
          type: 'component',
          name: 'API',
          pause: {
            message: 'Maintenance window',
            type: 'maintenance',
          },
        },
      },
    });
    const content = wrapper.findComponent(CTooltipStub).props('content');

    expect(content).toContain('Pause is inherited from Component API');
    expect(content).toContain('Maintenance window');
    expect(content).toContain('<div>Pause is inherited from Component API</div>');
    expect(content).toContain(
      '<div class="c-extra-details__message pre-line">Maintenance window</div>',
    );
  });

  it('matches equal multiline comments and preserves their line breaks', () => {
    const wrapper = factory({
      propsData: {
        pause: { m: 'First line\nSecond line', pause_type: 'maintenance' },
        entity: {
          type: 'service',
          name: 'API',
          pause: { message: 'First line\nSecond line', type: 'maintenance' },
        },
      },
    });
    const content = wrapper.findComponent(CTooltipStub).props('content');

    expect(content).toContain('Pause is inherited from Service API');
    expect(content).toContain('First line\nSecond line');
  });

  it('does not show the inherited message when comments differ', () => {
    const wrapper = factory({
      propsData: {
        pause: {
          m: 'Alarm pause comment',
          pause_type: 'maintenance',
        },
        entity: {
          type: 'component',
          name: 'API',
          pause: {
            message: 'Entity pause comment',
            type: 'maintenance',
          },
        },
      },
    });
    const content = wrapper.findComponent(CTooltipStub).props('content');

    expect(content).toContain('Alarm pause comment');
    expect(content).not.toContain('Pause is inherited from');
  });

  it('shows the inherited message regardless of legacy pause type metadata', () => {
    const wrapper = factory({
      propsData: {
        pause: { m: 'Maintenance window', pause_type: 'maintenance' },
        entity: {
          type: 'component',
          name: 'API',
          pause: { message: 'Maintenance window', type: 'planned maintenance' },
        },
      },
    });
    const content = wrapper.findComponent(CTooltipStub).props('content');

    expect(content).toContain('Pause is inherited from Component API');
  });
});
