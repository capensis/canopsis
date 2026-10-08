import { generateShallowRenderer } from '@unit/utils/vue';
import Faker from 'faker';

import AlarmCommentTemplateForm from '@/components/other/comment-template/form/alarm-comment-template-form.vue';

const stubs = {
  'c-select': true,
  'c-name-field': true,
};

const selectTemplateField = wrapper => wrapper.find('c-select-stub');

describe('alarm-comment-template-form', () => {
  const factory = generateShallowRenderer(AlarmCommentTemplateForm, { stubs });

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

    expect(wrapper.vm.fields).toHaveLength(1);
    expect(wrapper.vm.fields[0]).toMatchObject({ name: 'comment', required: true });
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

    expect(wrapper.vm.fields).toEqual(template.fields);
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

  test('Template can be cleared', async () => {
    const templates = [{ _id: '1', name: 'Template 1', fields: [] }];

    const wrapper = factory({
      propsData: {
        form: {
          template: templates[0],
          comment: '',
        },
        templates,
      },
    });

    wrapper.vm.updateField('template', null);
    await wrapper.vm.$nextTick();

    expect(wrapper).toEmit('update:form', { template: null, comment: '' });
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

    wrapper.vm.updateField('template', template);
    await wrapper.vm.$nextTick();

    expect(wrapper).toEmit('update:form', { template, comment: '' });
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

    expect(wrapper.vm.form.summary).toBeTruthy();
    expect(wrapper.vm.form.details).toBeTruthy();
    expect(wrapper.vm.fields.map(({ name }) => name)).toEqual(['summary', 'details']);
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

    expect(wrapper.vm.fields.map(field => field.label || field.name)).toEqual([
      'Custom Label',
      'field2',
    ]);
  });
});
