import Faker from 'faker';

import { flushPromises, generateRenderer } from '@unit/utils/vue';
import { mockXMLHttpRequest } from '@unit/utils/mock-hooks';
import { createMockedStoreModules } from '@unit/utils/store';

import { API_HOST, API_ROUTES } from '@/config';

import TextEditor from '@/components/common/text-editor/text-editor.vue';

const stubs = {
  'variables-menu': true,
};

const selectEditor = wrapper => wrapper.find('.jodit-wysiwyg');
const selectEditorImageControl = wrapper => wrapper.find('.jodit-toolbar-button_image button');
const selectEditorDragAndDropFileInput = wrapper => wrapper.find('.jodit-drag-and-drop__file-box input');
const selectEditorImageControlTabs = wrapper => wrapper.findAll('.jodit-tabs__buttons button');
const selectEditorImageControlUrlTab = wrapper => selectEditorImageControlTabs(wrapper).at(1);
const selectEditorTabsWrapper = wrapper => wrapper.find('.jodit-popup__content');
const selectEditorImageUrlInput = wrapper => selectEditorTabsWrapper(wrapper)
  .find('input[name="url"]');
const selectEditorImageTextInput = wrapper => selectEditorTabsWrapper(wrapper)
  .find('input[name="text"]');
const selectEditorImageInsetButton = wrapper => selectEditorTabsWrapper(wrapper)
  .findAll('button');
const selectVariablesMenu = wrapper => wrapper.find('variables-menu-stub');
const selectVariablesButton = wrapper => wrapper.find('.jodit-ui-group__variables button');

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

  const store = createMockedStoreModules([{
    name: 'templateVars',
    getters: {
      items: { test: 'test', testObj: { testObjField: 'value' } },
    },
  }]);

  test('Value changed after change props', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    const newValue = `<div>${Faker.lorem.words()}</div>`;

    await wrapper.setProps({ value: newValue });

    expect(wrapper).toEmitInput(newValue);
  });

  test('Value changed after trigger editor', async () => {
    const wrapper = snapshotFactory();

    await flushPromises();

    const editor = selectEditor(wrapper);

    const newValue = `<div>${Faker.lorem.words()}</div>`;

    editor.element.innerHTML = newValue;
    editor.trigger('mousedown');

    expect(wrapper).toEmitInput(newValue);
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
    range.setEnd(editor.element.firstChild, 1);
    selection.removeAllRanges();
    selection.addRange(range);
    const variablesMenu = selectVariablesMenu(wrapper);

    variablesMenu.triggerCustomEvent('input', variable);

    expect(wrapper).toEmitInput(`<p>${variable}</p>`);
    expect(focusSpy).toHaveBeenCalled();
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

    const variable = `{{ ${Faker.lorem.word()} }}`;

    const editor = selectEditor(wrapper);
    const range = document.createRange();
    const selection = window.getSelection();

    range.setStart(editor.element.firstChild.firstChild, 15);
    selection.removeAllRanges();
    selection.addRange(range);

    const variablesMenu = selectVariablesMenu(wrapper);

    variablesMenu.triggerCustomEvent('input', variable);

    await flushPromises();

    expect(wrapper.vm.editor.value).toEqual(`<p>Variable: ${variable}</p>`);
    expect(focusSpy).toHaveBeenCalled();
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

    const variable = `{{ ${Faker.lorem.word()} }}`;

    const editor = selectEditor(wrapper);
    const range = document.createRange();
    const selection = window.getSelection();

    range.setStart(editor.element.firstChild.firstChild, initialValue.indexOf('{{'));
    range.setEnd(editor.element.firstChild.firstChild, initialValue.indexOf('}}') + 2);
    selection.removeAllRanges();
    selection.addRange(range);

    const variablesMenu = selectVariablesMenu(wrapper);

    variablesMenu.triggerCustomEvent('input', variable);

    expect(wrapper.vm.editor.value).toEqual(`<p>Variable: ${variable}</p>`);
    expect(focusSpy).toHaveBeenCalled();
  });

  test('Image uploaded after trigger image control', async () => {
    jest.useFakeTimers('legacy');
    const focusSpy = jest.spyOn(window, 'focus').mockImplementation(() => {});

    const [file] = files;

    const wrapper = snapshotFactory();

    await flushPromises();

    const editorImageControl = selectEditorImageControl(wrapper);

    await editorImageControl.trigger('click');
    const editorDragAndDropFileInput = selectEditorDragAndDropFileInput(wrapper);
    Object.defineProperty(editorDragAndDropFileInput.element, 'files', {
      value: [file],
    });
    editorDragAndDropFileInput.trigger('change');

    await flushPromises();

    expect(XMLHttpRequest.open).toHaveBeenCalledWith('post', `${API_HOST}${API_ROUTES.file}?public=false`, true);

    await flushPromises();
    jest.runAllTimers();

    expect(XMLHttpRequest.send).toHaveBeenCalledWith(expect.any(FormData));

    const [fileResponse] = filesResponse;

    XMLHttpRequest.responseText = JSON.stringify([fileResponse]);
    XMLHttpRequest.status = 200;

    XMLHttpRequest.onload();

    await flushPromises();

    expect(focusSpy).toHaveBeenCalled();

    expect(wrapper).toEmitInput(
      `<p><img src="${API_HOST}${API_ROUTES.file}/${fileResponse._id}" width="300"></p><p><br></p>`,
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

    await editorImageControl.trigger('click');

    const editorImageControlUrlTab = selectEditorImageControlUrlTab(wrapper);
    editorImageControlUrlTab.trigger('click');

    const editorImageUrlInput = selectEditorImageUrlInput(wrapper);
    const editorImageTextInput = selectEditorImageTextInput(wrapper);
    const editorImageInsetButton = selectEditorImageInsetButton(wrapper);

    editorImageUrlInput.setValue(url);
    editorImageTextInput.setValue(text);
    editorImageInsetButton.trigger('click');

    expect(focusSpy).toHaveBeenCalled();

    expect(wrapper).toEmitInput(
      `<p><img src="${url}" alt="${text}" width="300"></p><p><br></p>`,
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

    await editorImageControl.trigger('click');

    const editorDragAndDropFileInput = selectEditorDragAndDropFileInput(wrapper);
    Object.defineProperty(editorDragAndDropFileInput.element, 'files', {
      value: [fileWithMaxSize],
    });
    editorDragAndDropFileInput.trigger('change');

    await flushPromises();

    expect(XMLHttpRequest.open).not.toHaveBeenCalled();
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

    expect(XMLHttpRequest.open).toHaveBeenCalledWith('post', `${API_HOST}${API_ROUTES.file}?public=false`, true);

    await flushPromises();
    jest.runAllTimers();

    expect(XMLHttpRequest.send).toHaveBeenCalledWith(expect.any(FormData));

    XMLHttpRequest.responseText = JSON.stringify(filesResponse);
    XMLHttpRequest.status = 200;

    XMLHttpRequest.onload();

    await flushPromises();

    expect(focusSpy).toHaveBeenCalled();

    const firstFile = filesResponse[0];
    const secondFile = filesResponse[1];

    const firstEmitData = `<p><img src="${API_HOST}${API_ROUTES.file}/${firstFile._id}" width="300"></p><p><br></p>`;

    const secondEmitData = `<p><img src="${API_HOST}${API_ROUTES.file}/${firstFile._id}" width="300"><a href="${API_HOST}${API_ROUTES.file}/${secondFile._id}" target="_blank">${secondFile.filename}</a></p><p><br></p>`;

    expect(wrapper).toEmitInput(
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
      store,
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
