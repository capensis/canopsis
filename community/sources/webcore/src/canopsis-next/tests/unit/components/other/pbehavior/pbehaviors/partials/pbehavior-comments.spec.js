import { range } from 'lodash';

import { createVueInstance, generateRenderer } from '@unit/utils/vue';

import PbehaviorComments from '@/components/other/pbehavior/pbehaviors/partials/pbehavior-comments.vue';

const localVue = createVueInstance();

describe('pbehavior-comments', () => {
  const totalItems = 5;

  const pbehaviorComments = range(totalItems).map(index => ({
    _id: `id-pbehavior-comment-${index}`,
    author: {
      name: `author-pbehavior-comment-${index}`,
    },
    message: `message-pbehavior-comment-${index}`,
  }));

  const snapshotFactory = generateRenderer(PbehaviorComments, { localVue });

  test('Renders `pbehavior-comments` without comments', () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `pbehavior-comments` with comments', () => {
    const wrapper = snapshotFactory({
      propsData: {
        comments: pbehaviorComments,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
