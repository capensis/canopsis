import Faker from 'faker';

import { createVueInstance, generateRenderer, generateShallowRenderer } from '@unit/utils/vue';
import { createTextareaInputStub } from '@unit/stubs/input';

import PbehaviorCommentField from '@/components/other/pbehavior/pbehaviors/fields/pbehavior-comment-field.vue';

const localVue = createVueInstance();

const stubs = {
  'v-textarea': createTextareaInputStub('v-textarea'),
};

const selectMessageField = wrapper => wrapper.find('.v-textarea');
const selectRemoveButton = wrapper => wrapper.find('v-btn-stub');

describe('pbehavior-comment-field', () => {
  const factory = generateShallowRenderer(PbehaviorCommentField, {
    localVue,
    stubs,
  });
  const snapshotFactory = generateRenderer(PbehaviorCommentField, {
    localVue,
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

    selectMessageField(wrapper).vm.$emit('input', newMessage);

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

    selectRemoveButton(wrapper).vm.$emit('click', newMessage);

    expect(wrapper).toEmit('remove');
  });

  test('Renders `pbehavior-comment-field` with default props', () => {
    const wrapper = snapshotFactory({
      propsData: {
        comment: {},
      },
    });

    expect(wrapper.element).toMatchSnapshot();
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

    expect(wrapper.element).toMatchSnapshot();
  });
});
