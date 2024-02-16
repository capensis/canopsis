import flushPromises from 'flush-promises';

import { generateRenderer } from '@unit/utils/vue';

import CRuntimeTemplate from '@/components/common/runtime-template/c-runtime-template.vue';
import CCompiledTemplate from '@/components/common/runtime-template/c-compiled-template.vue';
import AlarmsExpandPanelMoreInfos from '@/components/widgets/alarm/expand-panel/alarms-expand-panel-more-infos.vue';

const stubs = {
  'c-compiled-template': CCompiledTemplate,
  'c-runtime-template': CRuntimeTemplate,
};

describe('alarms-expand-panel-more-infos', () => {
  const snapshotFactory = generateRenderer(AlarmsExpandPanelMoreInfos, {
    stubs,
    parentComponent: {
      provide: {
        $system: {},
      },
    },
  });

  it('Renders `alarms-expand-panel-more-infos` without template', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
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

    expect(wrapper).toMatchSnapshot();
  });
});
