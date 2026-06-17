import { generateRenderer } from '@unit/utils/vue';

import CList from '@/components/common/list/c-list.vue';
import VariablesList from '@/components/common/text-editor/variables-list.vue';

describe('variables-list', () => {
  const snapshotFactory = generateRenderer(VariablesList, {
    attachTo: document.body,
    components: {
      CList,
    },
    stubs: {
      'v-list-item-mask': true,
    },
  });

  test('Renders `variables-list` with default props', () => {
    const wrapper = snapshotFactory();

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `variables-list` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: 'entity._id',
        zIndex: 2,
        childrenKey: 'variables',
        items: [
          {
            value: 'entity',
            text: 'Entity',
            variables: [
              {
                value: 'id',
                text: 'Id',
              },
            ],
          },
        ],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
