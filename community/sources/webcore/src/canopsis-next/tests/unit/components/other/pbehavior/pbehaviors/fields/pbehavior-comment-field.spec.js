import Faker from 'faker';

import { generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createTextareaInputStub } from '@unit/stubs/input';

import PbehaviorCommentField from '@/components/other/pbehavior/pbehaviors/fields/pbehavior-comment-field.vue';

const stubs = {
  'v-textarea': createTextareaInputStub('v-textarea'),
  'c-action-btn': true,
};

const snapshotStubs = {
  'c-action-btn': true,
};

const selectMessageField = wrapper => wrapper.find('.v-textarea');
const selectRemoveButton = wrapper => wrapper.find('c-action-btn-stub');

describe('pbehavior-comment-field', () => {
  const factory = generateShallowRenderer(PbehaviorCommentField, { stubs });
  const snapshotFactory = generateRenderer(PbehaviorCommentField, {
    stubs: snapshotStubs,
  });

  test('Message changed after trigger textarea', () => {
    const comment = {
      key: Faker.datatype.string(),
      message: Faker.datatype.string(),
    };
    const wrapper = factory({
      propsData: { comment },
    });

    const newMessage = Faker.datatype.string();

    selectMessageField(wrapper).triggerCustomEvent('input', newMessage);

    expect(wrapper).toEmit('input', {
      ...comment,
      message: newMessage,
    });
  });

  test('Remove event emitted after trigger remove button', () => {
    const wrapper = factory({
      propsData: { comment: {} },
    });

    const newMessage = Faker.datatype.string();

    selectRemoveButton(wrapper).triggerCustomEvent('click', newMessage);

    expect(wrapper).toEmit('remove');
  });

  test('Renders `pbehavior-comment-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        comment: {},
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `pbehavior-comment-field` with custom props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        comment: {
          message: 'comment-message',
        },
        max: 100,
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
