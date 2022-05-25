import flushPromises from 'flush-promises';
import Faker from 'faker';

import { mount, createVueInstance } from '@unit/utils/vue';
import { mockXMLHttpRequest } from '@unit/utils/mock-hooks';
import { API_HOST, API_ROUTES } from '@/config';

import TextEditor from '@/components/common/text-editor/text-editor.vue';

const localVue = createVueInstance();

const snapshotFactory = (options = {}) => mount(TextEditor, {
  localVue,
  attachTo: document.body,

  ...options,
});

const selectEditor = wrapper => wrapper.find('.jodit_wysiwyg');
const selectEditorImageControl = wrapper => wrapper.find('.jodit_toolbar_btn-image');
const selectEditorDragAndDropFileInput = wrapper => wrapper.find('.jodit_draganddrop_file_box input');
const selectEditorImageControlTabs = wrapper => selectEditorImageControl(wrapper)
  .findAll('.jodit_tabs_buttons a');
const selectEditorImageControlUrlTab = wrapper => selectEditorImageControlTabs(wrapper).at(1);
const selectEditorTabsWrapper = wrapper => selectEditorImageControl(wrapper)
  .find('.jodit_tabs_wrapper');
const selectEditorImageUrlInput = wrapper => selectEditorTabsWrapper(wrapper)
  .find('input[name="url"]');
const selectEditorImageTextInput = wrapper => selectEditorTabsWrapper(wrapper)
  .find('input[name="text"]');
const selectEditorImageInsetButton = wrapper => selectEditorTabsWrapper(wrapper)
  .findAll('button');

describe('text-editor', () => {
  const XMLHttpRequest = mockXMLHttpRequest();

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

  test('Image uploaded after trigger image control', async () => {
    const focusSpy = jest.spyOn(window, 'focus').mockImplementation(() => {});

    const filename = 'file.png';
    const mediatype = 'image/png';
    const file = new File([new ArrayBuffer(1)], filename, { type: mediatype });

    const wrapper = snapshotFactory();

    await flushPromises();

    const editorImageControl = selectEditorImageControl(wrapper);

    await editorImageControl.trigger('mousedown');

    const editorDragAndDropFileInput = selectEditorDragAndDropFileInput(wrapper);
    Object.defineProperty(editorDragAndDropFileInput.element, 'files', {
      value: [file],
    });
    editorDragAndDropFileInput.trigger('change');

    await flushPromises();

    expect(XMLHttpRequest.open).toBeCalledWith('POST', `${API_HOST}${API_ROUTES.file}?public=false`, true);

    await flushPromises();

    expect(XMLHttpRequest.send).toBeCalledWith(expect.any(FormData));

    const fileResponse = {
      _id: '9f6e9bbe-c021-4ede-94bd-b42dfa7d22f7',
      filename,
      mediatype,
      created: 1653398538,
    };

    XMLHttpRequest.responseText = JSON.stringify([fileResponse]);
    XMLHttpRequest.status = 200;

    XMLHttpRequest.onload();

    await flushPromises();

    expect(focusSpy).toBeCalled();

    expect(wrapper).toEmit(
      'input',
      `<img src="${API_HOST}${API_ROUTES.file}/${fileResponse._id}" style="width: 300px;">`,
    );
  });

  test('Image as url uploaded after trigger image control', async () => {
    const focusSpy = jest.spyOn(window, 'focus').mockImplementation(() => {});
    const url = Faker.lorem.word();
    const text = Faker.lorem.word();
    const wrapper = snapshotFactory();

    await flushPromises();

    const editorImageControl = selectEditorImageControl(wrapper);

    await editorImageControl.trigger('mousedown');

    const editorImageControlUrlTab = selectEditorImageControlUrlTab(wrapper);
    editorImageControlUrlTab.trigger('click');

    const editorImageUrlInput = selectEditorImageUrlInput(wrapper);
    const editorImageTextInput = selectEditorImageTextInput(wrapper);
    const editorImageInsetButton = selectEditorImageInsetButton(wrapper);

    editorImageUrlInput.setValue(url);
    editorImageTextInput.setValue(text);
    editorImageInsetButton.trigger('click');

    expect(focusSpy).toBeCalled();

    expect(wrapper).toEmit(
      'input',
      `<img src="${url}" alt="${text}" style="width: 300px;">`,
    );
  });

  test('Image not uploaded after trigger image control with large file', async () => {
    const filename = 'file.png';
    const mediatype = 'image/png';
    const fileWithMaxSize = new File([new ArrayBuffer(2)], filename, { type: mediatype });

    const wrapper = snapshotFactory({
      propsData: {
        maxFileSize: 1,
      },
    });

    await flushPromises();

    const editorImageControl = selectEditorImageControl(wrapper);

    await editorImageControl.trigger('mousedown');

    const editorDragAndDropFileInput = selectEditorDragAndDropFileInput(wrapper);
    Object.defineProperty(editorDragAndDropFileInput.element, 'files', {
      value: [fileWithMaxSize],
    });
    editorDragAndDropFileInput.trigger('change');

    await flushPromises();

    expect(XMLHttpRequest.open).not.toBeCalled();
  });

  test('Renders `text-editor` with default props', async () => {
    jest.useFakeTimers('legacy');
    const wrapper = snapshotFactory();

    await flushPromises();
    jest.runAllTimers();

    expect(wrapper.element).toMatchSnapshot();
    jest.useRealTimers();
  });

  test('Renders `text-editor` with files', async () => {
    jest.useFakeTimers('legacy');
    const wrapper = snapshotFactory({
      value: `<img src="${API_HOST}${API_ROUTES.file}/123" style="width: 300px;">`,
    });

    await flushPromises();
    jest.runAllTimers();

    expect(wrapper.element).toMatchSnapshot();
    jest.useRealTimers();
  });

  test('Renders `text-editor` with custom props', async () => {
    jest.useFakeTimers('legacy');
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

    await flushPromises();
    jest.runAllTimers();

    expect(wrapper.element).toMatchSnapshot();
    jest.useRealTimers();
  });
});
