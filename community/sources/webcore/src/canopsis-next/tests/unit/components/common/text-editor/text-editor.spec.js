import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, createVueInstance } from '@unit/utils/vue';

import TextEditor from '@/components/common/text-editor/text-editor.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(TextEditor, {
  localVue,
  attachTo: document.body,

  ...options,
});

const selectEditor = wrapper => wrapper.find('.jodit_wysiwyg');

describe('text-editor', () => {
  beforeAll(() => {
    Object.defineProperty(HTMLElement.prototype, 'innerText', {
      set() {},
      get() {
        return this.textContent;
      },
    });
  });

  test('Value changed after change props', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    const newValue = `<div>${Faker.lorem.words()}</div>`;

    await wrapper.setProps({ value: newValue });

    expect(wrapper).toEmit('input', newValue);
  });

  test('Value changed after trigger editor', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    const editor = selectEditor(wrapper);

    const newValue = `<div>${Faker.lorem.words()}</div>`;

    editor.element.innerHTML = newValue;
    editor.trigger('mousedown');

    expect(wrapper).toEmit('input', newValue);
  });

  test('Renders `text-editor` with default props', async () => {
    const wrapper = snapshotFactory();

    expect(wrapper.element).toMatchSnapshot();
  });

  test('Renders `text-editor` with custom props', async () => {
    const wrapper = snapshotFactory({
      propsData: {
        value: '<div><p>Paragraph</p></div>',
        label: 'Text editor label',
        buttons: [{}],
        public: true,
        extraButtons: [{}],
        config: {},
        errorMessages: ['Error'],
        maxFileSize: 1,
      },
    });

    expect(wrapper.element).toMatchSnapshot();
  });
});
