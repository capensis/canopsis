import Faker from 'faker';

import { generateRenderer, generateShallowRenderer, flushPromises } from '@unit/utils/vue';
import { createInputStub, createSelectInputStub } from '@unit/stubs/input';

import AlarmCommentTemplateForm from '@/components/other/comment-template/form/alarm-comment-template-form.vue';

const stubs = {
  'c-select-field': createSelectInputStub('c-select-field'),
  'c-name-field': createInputStub('c-name-field'),
};

const snapshotStubs = {
  'c-select-field': true,
  'c-name-field': true,
};

const selectTemplateField = wrapper => wrapper.find('.c-select-field');
const selectNameFields = wrapper => wrapper.findAll('.c-name-field');

describe('alarm-comment-template-form', () => {
  const factory = generateShallowRenderer(AlarmCommentTemplateForm, { stubs });
  const snapshotFactory = generateRenderer(AlarmCommentTemplateForm, { stubs: snapshotStubs });

  test('Renders default comment field when no template is selected', () => {
    const wrapper = factory({
      propsData: {
        form: {
          template: null,
          comment: '',
        },
        templates: [],
      },
    });

    const nameFields = selectNameFields(wrapper);

    expect(nameFields).toHaveLength(1);
    expect(nameFields.at(0).attributes('name')).toBe('comment');
    expect(nameFields.at(0).attributes('required')).toBe('required');
  });

  test('Renders template fields when template is selected', () => {
    const template = {
      _id: '1',
      name: 'Template 1',
      fields: [
        { name: 'field1', label: 'Field 1', required: true },
        { name: 'field2', label: 'Field 2', required: false },
        { name: 'field3', label: 'Field 3', required: true },
      ],
    };

    const wrapper = factory({
      propsData: {
        form: {
          template,
          field1: '',
          field2: '',
          field3: '',
        },
        templates: [template],
      },
    });

    const nameFields = selectNameFields(wrapper);

    expect(nameFields).toHaveLength(3);
    expect(nameFields.at(0).attributes('name')).toBe('field1');
    expect(nameFields.at(0).attributes('required')).toBe('required');
    expect(nameFields.at(1).attributes('name')).toBe('field2');
    expect(nameFields.at(1).attributes('required')).toBeUndefined();
    expect(nameFields.at(2).attributes('name')).toBe('field3');
    expect(nameFields.at(2).attributes('required')).toBe('required');
  });

  test('Does not render template select when templates array is empty', () => {
    const wrapper = factory({
      propsData: {
        form: {
          template: null,
          comment: '',
        },
        templates: [],
      },
    });

    const templateField = selectTemplateField(wrapper);

    expect(templateField.exists()).toBe(false);
  });

  test('Clears errors when template is changed', async () => {
    const templates = [
      { _id: '1', name: 'Template 1', fields: [] },
    ];

    const wrapper = factory({
      propsData: {
        form: {
          template: null,
          comment: '',
        },
        templates,
      },
    });

    const validator = wrapper.getValidator();
    const clearSpy = jest.spyOn(validator.errors, 'clear');

    const templateField = selectTemplateField(wrapper);
    templateField.triggerCustomEvent('input', templates[0]);

    await flushPromises();

    expect(clearSpy).toHaveBeenCalled();
  });

  test('Template select field emits input with correct template object', async () => {
    const template = {
      _id: '1',
      name: 'Bug Report Template',
      fields: [
        { name: 'summary', label: 'Summary', required: true },
        { name: 'steps', label: 'Steps to Reproduce', required: true },
        { name: 'expected', label: 'Expected Result', required: false },
      ],
    };
    const templates = [template];

    const wrapper = factory({
      propsData: {
        form: {
          template: null,
          comment: '',
        },
        templates,
      },
    });

    const templateField = selectTemplateField(wrapper);

    expect(templateField.vm.value).toBeNull();

    templateField.triggerCustomEvent('input', template);

    await flushPromises();

    expect(wrapper).toEmitInput({ template, comment: '' });
  });

  test('Fields are bound correctly with v-field directive', async () => {
    const template = {
      _id: '1',
      name: 'Template 1',
      fields: [
        { name: 'summary', label: 'Summary', required: true },
        { name: 'details', label: 'Details', required: false },
      ],
    };

    const wrapper = factory({
      propsData: {
        form: {
          template,
          summary: Faker.datatype.string(),
          details: Faker.datatype.string(),
        },
        templates: [template],
      },
    });

    const nameFields = selectNameFields(wrapper);

    expect(nameFields.at(0).vm.value).toBe(wrapper.vm.form.summary);
    expect(nameFields.at(1).vm.value).toBe(wrapper.vm.form.details);
  });

  test('Field labels use field.label or fallback to field.name', () => {
    const template = {
      _id: '1',
      name: 'Template 1',
      fields: [
        { name: 'field1', label: 'Custom Label', required: true },
        { name: 'field2', required: false },
      ],
    };

    const wrapper = factory({
      propsData: {
        form: {
          template,
          field1: '',
          field2: '',
        },
        templates: [template],
      },
    });

    const nameFields = selectNameFields(wrapper);

    expect(nameFields.at(0).attributes('label')).toBe('Custom Label');
    expect(nameFields.at(1).attributes('label')).toBe('field2');
  });

  test('Renders `alarm-comment-template-form` without template', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          template: null,
          comment: '',
        },
        templates: [],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });

  test('Renders `alarm-comment-template-form` with template', () => {
    const wrapper = snapshotFactory({
      propsData: {
        form: {
          template: {
            _id: '1',
            name: 'Template 1',
            fields: [
              { name: 'field1', label: 'Field 1', required: true },
              { name: 'field2', label: 'Field 2', required: false },
            ],
          },
          field1: '',
          field2: '',
        },
        templates: [
          {
            _id: '1',
            name: 'Template 1',
            fields: [
              { name: 'field1', label: 'Field 1', required: true },
              { name: 'field2', label: 'Field 2', required: false },
            ],
          },
        ],
      },
    });

    expect(wrapper).toMatchSnapshot();
  });
});
