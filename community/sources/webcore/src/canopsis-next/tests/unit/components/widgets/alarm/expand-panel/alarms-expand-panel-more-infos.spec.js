import flushPromises from 'flush-promises';

import { mount, createVueInstance } from '@unit/utils/vue';

import AlarmsExpandPanelMoreInfos from '@/components/widgets/alarm/expand-panel/alarms-expand-panel-more-infos.vue';

const localVue = createVueInstance();

const stubs = {
  'v-runtime-template': true,
};

const snapshotFactory = (options = {}) => mount(AlarmsExpandPanelMoreInfos, {
  localVue,
  stubs,

  ...options,
});

describe('alarms-expand-panel-more-infos', () => {
  it('Renders `alarms-expand-panel-more-infos` without template', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `alarms-expand-panel-more-infos` with template', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        template: '<div><p>Template</p><p>{{ alarm.name }}</p><p>{{ entity.name }}</p></div>',
        alarm: {
          name: 'more-infos-alarm',
          entity: {
            name: 'more-infos-entity',
          },
        },
      },
    });

    await flushPromises();

    expect(wrapper.element).toMatchSnapshot();
  });
});
