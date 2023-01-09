import { range } from 'lodash';

import { createVueInstance, generateRenderer } from '@unit/utils/vue';

import PbehaviorExceptionList from '@/components/other/pbehavior/pbehaviors/partials/pbehavior-exception-list.vue';

const localVue = createVueInstance();

describe('pbehavior-exception-list', () => {
  const totalItems = 5;
  const exceptions = range(totalItems).map(index => ({
    name: `exception-${index}`,
  }));

  const snapshotFactory = generateRenderer(PbehaviorExceptionList, { localVue });

  test('Renders `pbehavior-exception-list` without exceptions', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pbehavior-exception-list` with exceptions', () => {
    const wrapper = snapshotFactory({
      propsData: {
        exceptions,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
