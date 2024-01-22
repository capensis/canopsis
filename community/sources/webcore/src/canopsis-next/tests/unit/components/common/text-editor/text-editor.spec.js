import flushPromises from 'flush-promises';
import Faker from 'faker';

import { generateRenderer } from '@unit/utils/vue';
import { mockXMLHttpRequest } from '@unit/utils/mock-hooks';
import { API_HOST, API_ROUTES } from '@/config';

import TextEditor from '@/components/common/text-editor/text-editor.vue';

const stubs = {
  'variables-menu': true,
};

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
const selectVariablesMenu = wrapper => wrapper.find('variables-menu-stub');
const selectVariablesButton = wrapper => wrapper.find('.text-editor__variables-button');

describe('text-editor', () => {
  const XMLHttpRequest = mockXMLHttpRequest();
  const filesMeta = [
    {
      filename: 'file.png',
      mediatype: 'image/png',
    },
    {
      filename: 'file.doc',
      mediatype: 'other/doc',
    },
  ];
  const files = filesMeta.map(({ filename, mediatype }) => new File(
    [new ArrayBuffer(1)],
    filename,
    { type: mediatype },
  ));

  const filesResponse = filesMeta.map(({ filename, mediatype }) => ({
    _id: filename,
    filename,
    mediatype,
    created: 1653398538,
  }));

  const snapshotFactory = generateRenderer(TextEditor, {
    attachTo: document.body,
    stubs,
  }, { noDestroy: true });

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

  test('Value changed after trigger variables', async () => {
    const focusSpy = jest.spyOn(window, 'focus').mockImplementation();
    const initialValue = 'Text';
    const wrapper = snapshotFactory({
      propsData: {
        value: initialValue,
        variables: [{ value: 'variable' }],
      },
    });

    await flushPromises();

    const variable = Faker.lorem.word();
    const editor = selectEditor(wrapper);
    const range = document.createRange();
    const selection = window.getSelection();

    range.setStart(editor.element.firstChild, 0);
    range.setEnd(editor.element.firstChild, initialValue.length);
    selection.removeAllRanges();
    selection.addRange(range);

    const variablesMenu = selectVariablesMenu(wrapper);

    variablesMenu.vm.$emit('input', variable);

    expect(wrapper).toEmit('input', `{{ ${variable} }}`);
    expect(focusSpy).toBeCalled();
  });

  test('Menu showed after trigger variables button', async () => {
    const initialValue = 'Variable: {{ variable }}';
    const wrapper = snapshotFactory({
      propsData: {
        value: initialValue,
        variables: [{ value: 'variable' }],
      },
    });

    await flushPromises();

    const variablesButton = selectVariablesButton(wrapper);
    jest.spyOn(variablesButton.element, 'getBoundingClientRect').mockImplementation(() => ({
      top: 100,
      left: 100,
      height: 88,
    }));
    variablesButton.trigger('click');

    expect(wrapper.vm.variablesShown).toBeTruthy();
    expect(wrapper.vm.variablesMenuPosition).toEqual({
      y: 188,
      x: 100,
    });
  });

  test('Value changed after trigger variables with caret in variable', async () => {
    const focusSpy = jest.spyOn(window, 'focus').mockImplementation();
    const initialValue = 'Variable: {{ variable }}';
    const wrapper = snapshotFactory({
      propsData: {
        value: initialValue,
        variables: [{ value: 'variable' }],
      },
    });

    await flushPromises();

    const variable = Faker.lorem.word();

    const editor = selectEditor(wrapper);
    const range = document.createRange();
    const selection = window.getSelection();

    range.setStart(editor.element.firstChild, 15);
    selection.removeAllRanges();
    selection.addRange(range);

    const variablesMenu = selectVariablesMenu(wrapper);

    variablesMenu.vm.$emit('input', variable);

    expect(wrapper).toEmit('input', `Variable: {{ ${variable} }}`);
    expect(focusSpy).toBeCalled();
  });

  test('Value changed after trigger variables with selected variable', async () => {
    const focusSpy = jest.spyOn(window, 'focus').mockImplementation();
    const initialValue = 'Variable: {{ variable }}';
    const wrapper = snapshotFactory({
      propsData: {
        value: initialValue,
        variables: [{ value: 'variable' }],
      },
    });

    await flushPromises();

    const variable = Faker.lorem.word();

    const editor = selectEditor(wrapper);
    const range = document.createRange();
    const selection = window.getSelection();

    range.setStart(editor.element.firstChild, initialValue.indexOf('{{'));
    range.setEnd(editor.element.firstChild, initialValue.indexOf('}}') + 2);
    selection.removeAllRanges();
    selection.addRange(range);

    const variablesMenu = selectVariablesMenu(wrapper);

    variablesMenu.vm.$emit('input', variable);

    expect(wrapper).toEmit('input', `Variable: {{ ${variable} }}`);
    expect(focusSpy).toBeCalled();
  });

  test('Image uploaded after trigger image control', async () => {
    jest.useFakeTimers('legacy');
    const focusSpy = jest.spyOn(window, 'focus').mockImplementation(() => {});

    const [file] = files;

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
    jest.runAllTimers();

    expect(XMLHttpRequest.send).toBeCalledWith(expect.any(FormData));

    const [fileResponse] = filesResponse;

    XMLHttpRequest.responseText = JSON.stringify([fileResponse]);
    XMLHttpRequest.status = 200;

    XMLHttpRequest.onload();

    await flushPromises();

    expect(focusSpy).toBeCalled();

    expect(wrapper).toEmit(
      'input',
      `<img src="${API_HOST}${API_ROUTES.file}/${fileResponse._id}" style="width: 300px;">`,
    );
    jest.useRealTimers();
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

  test('Image not uploaded after trigger paste event', async () => {
    jest.useFakeTimers('legacy');
    const focusSpy = jest.spyOn(window, 'focus').mockImplementation(() => {});

    const wrapper = snapshotFactory();

    await flushPromises();

    const editor = selectEditor(wrapper);

    await editor.trigger('paste', {
      clipboardData: {
        files,
        getData: () => {},
      },
    });

    await flushPromises();

    expect(XMLHttpRequest.open).toBeCalledWith('POST', `${API_HOST}${API_ROUTES.file}?public=false`, true);

    await flushPromises();
    jest.runAllTimers();

    expect(XMLHttpRequest.send).toBeCalledWith(expect.any(FormData));

    XMLHttpRequest.responseText = JSON.stringify(filesResponse);
    XMLHttpRequest.status = 200;

    XMLHttpRequest.onload();

    await flushPromises();

    expect(focusSpy).toBeCalled();

    const firstFile = filesResponse[0];
    const secondFile = filesResponse[1];

    const firstEmitData = `<img src="${API_HOST}${API_ROUTES.file}/${firstFile._id}" style="width: 300px;">`;

    const secondEmitData = `${firstEmitData}<a href="${API_HOST}${API_ROUTES.file}/${secondFile._id}" target="_blank">${secondFile.filename}</a>`;

    expect(wrapper).toEmit(
      'input',
      firstEmitData,
      secondEmitData,
    );
    jest.useRealTimers();
  });

  test('Renders `text-editor` with default props', async () => {
    jest.useFakeTimers();
    const wrapper = snapshotFactory();

    await flushPromises();
    jest.runAllTimers();

    expect(wrapper).toMatchSnapshot();
    jest.useRealTimers();
  });

  test('Renders `text-editor` with files', async () => {
    jest.useFakeTimers();
    const wrapper = snapshotFactory({
      value: `<img src="${API_HOST}${API_ROUTES.file}/123" style="width: 300px;">`,
    });

    await flushPromises();
    jest.runAllTimers();

    expect(wrapper).toMatchSnapshot();
    jest.useRealTimers();
  });

  test('Renders `text-editor` with custom props', async () => {
    jest.useFakeTimers();
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

    expect(wrapper).toMatchSnapshot();
    jest.useRealTimers();
  });

  test('Renders `text-editor` with variables', async () => {
    jest.useFakeTimers();
    const wrapper = snapshotFactory({
      propsData: {
        value: '{{ test.test }}',
        variables: [{ value: 'test.test' }],
      },
    });

    await flushPromises();
    jest.runAllTimers();

    const variablesButton = selectVariablesButton(wrapper);
    jest.spyOn(variablesButton.element, 'getBoundingClientRect').mockImplementation(() => ({
      top: 101,
      left: 112,
      height: 88,
    }));
    variablesButton.trigger('click');

    await flushPromises();

    expect(wrapper).toMatchSnapshot();
    jest.useRealTimers();
  });
});
