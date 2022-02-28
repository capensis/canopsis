import flushPromises from 'flush-promises';

import { mount, createVueInstance } from '@unit/utils/vue';

import MoreInfos from '@/components/widgets/alarm/more-infos/more-infos.vue';

const localVue = createVueInstance();

const stubs = {
  'v-runtime-template': true,
};

const snapshotFactory = (options = {}) => mount(MoreInfos, {
  localVue,
  stubs,

  ...options,
});

describe('more-infos', () => {
  it('Renders `more-infos` without template', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  it('Renders `more-infos` with template', async () => {
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
